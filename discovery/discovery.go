package discovery

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/net/ethernet"
	"github.com/kandoo/beehive-netctrl/nom"
)

type LinkDiscovered struct {
}

const (
	nodeDict = "N"
)

type nodePortsAndLinks struct {
	N nom.Node
	P []nom.Port
	L []nom.Link
}

func (np *nodePortsAndLinks) hasPort(port nom.Port) bool {
	for _, p := range np.P {
		if p.ID == port.ID {
			return true
		}
	}
	return false
}

func (np *nodePortsAndLinks) removePort(port nom.Port) bool {
	for i, p := range np.P {
		if p.ID == port.ID {
			np.P = append(np.P[:i], np.P[i+1:]...)
			return true
		}
	}
	return false
}

func (np *nodePortsAndLinks) hasLinkFrom(from nom.UID) bool {
	for _, l := range np.L {
		if l.From == from {
			return true
		}
	}
	return false
}

func (np *nodePortsAndLinks) hasLink(link nom.Link) bool {
	id := link.UID()
	for _, l := range np.L {
		if l.UID() == id {
			return true
		}
	}
	return false
}

func (np *nodePortsAndLinks) removeLink(link nom.Link) bool {
	for i, l := range np.L {
		if l.From == link.From {
			np.L = append(np.L[:i], np.L[i+1:]...)
			return true
		}
	}
	return false
}

type nodeJoinedHandler struct{}

func (h *nodeJoinedHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	joined := msg.Data().(nom.NodeJoined)
	d := ctx.Dict(nodeDict)
	n := nom.Node(joined)
	k := string(n.UID())
	var np nodePortsAndLinks
	if err := d.GetGob(k, &np); err != nil {
		glog.Warningf("%v rejoins", n)
	}
	np.N = n
	// TODO(soheil): Add a flow entry to forward lldp packets to the controller.
	return d.PutGob(k, &np)
}

func (h *nodeJoinedHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return bh.MappedCells{
		{nodeDict, string(nom.Node(msg.Data().(nom.NodeJoined)).UID())},
	}
}

type nodeLeftHandler struct{}

func (h *nodeLeftHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	n := nom.Node(msg.Data().(nom.NodeLeft))
	d := ctx.Dict(nodeDict)
	k := string(n.UID())
	if _, err := d.Get(k); err != nil {
		return fmt.Errorf("%v is not joined", n)
	}
	d.Del(k)
	return nil
}

func (h *nodeLeftHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return bh.MappedCells{
		{nodeDict, string(nom.Node(msg.Data().(nom.NodeLeft)).UID())},
	}
}

type portUpdateHandler struct{}

func (h *portUpdateHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	p := nom.Port(msg.Data().(nom.PortUpdated))
	d := ctx.Dict(nodeDict)
	k := string(p.Node)
	var np nodePortsAndLinks
	if err := d.GetGob(k, &np); err != nil {
		glog.Warningf("%v added before its node", p)
		ctx.Snooze(1 * time.Second)
		return nil
	}

	if np.hasPort(p) {
		glog.Warningf("%v readded")
		np.removePort(p)
	}

	sendLLDPPacket(np.N, p, ctx)

	np.P = append(np.P, p)
	return d.PutGob(k, &np)
}

func (h *portUpdateHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return bh.MappedCells{
		{nodeDict, string(msg.Data().(nom.PortUpdated).Node)},
	}
}

type lldpTimeout struct{}

type timeoutHandler struct{}

func (h *timeoutHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	d := ctx.Dict(nodeDict)
	d.ForEach(func(k string, v []byte) {
		var np nodePortsAndLinks
		if err := nom.ObjGoDecode(&np, v); err != nil {
			glog.Errorf("Error in decoding value: %v", err)
			return
		}
		for _, p := range np.P {
			sendLLDPPacket(np.N, p, ctx)
		}
	})
	return nil
}

func (h *timeoutHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return bh.MappedCells{}
}

type pktInHandler struct{}

func (h *pktInHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	pin := msg.Data().(nom.PacketIn)
	e := ethernet.NewEthernetWithBuf([]byte(pin.Packet))
	if e.Type() != uint16(ethernet.ETH_T_LLDP) {
		return nil
	}

	_, port, err := decodeLLDP([]byte(pin.Packet))
	if err != nil {
		return err
	}

	d := ctx.Dict(nodeDict)
	k := string(pin.Node)
	var np nodePortsAndLinks
	if err := d.GetGob(k, &np); err != nil {
		return fmt.Errorf("Node %v not found", pin.Node)
	}

	l := nom.Link{
		ID:    nom.LinkID(port.UID()),
		From:  pin.InPort,
		To:    []nom.UID{port.UID()},
		State: nom.LinkStateUp,
	}
	ctx.Emit(NewLink(l))

	l = nom.Link{
		ID:    nom.LinkID(pin.InPort),
		From:  port.UID(),
		To:    []nom.UID{pin.InPort},
		State: nom.LinkStateUp,
	}
	ctx.Emit(NewLink(l))

	return nil
}

func (h *pktInHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return bh.MappedCells{
		{nodeDict, string(msg.Data().(nom.PacketIn).Node)},
	}
}

type NewLink nom.Link

type newLinkHandler struct{}

func (h *newLinkHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	l := nom.Link(msg.Data().(NewLink))
	n, _ := nom.ParsePortUID(l.From)
	d := ctx.Dict(nodeDict)
	k := string(n)
	var np nodePortsAndLinks
	if err := d.GetGob(k, &np); err != nil {
		return err
	}

	if np.hasLinkFrom(l.From) {
		if np.hasLink(l) {
			return nil
		}
		np.removeLink(l)
		ctx.Emit(nom.LinkRemoved(l))
	}

	glog.V(2).Infof("Link detected %v", l)
	ctx.Emit(nom.LinkAdded(l))
	np.L = append(np.L, l)
	return d.PutGob(k, &np)
}

func (h *newLinkHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	n, _ := nom.ParsePortUID(msg.Data().(NewLink).From)
	return bh.MappedCells{{nodeDict, string(n)}}
}

// RegisterDiscovery registers the handlers for topology discovery on the hive.
func RegisterDiscovery(h bh.Hive) {
	a := h.NewApp("discovery")
	a.Handle(nom.NodeJoined{}, &nodeJoinedHandler{})
	a.Handle(nom.NodeLeft{}, &nodeLeftHandler{})
	a.Handle(nom.PortUpdated{}, &portUpdateHandler{})
	// TODO(soheil): Handle PortRemoved.
	a.Handle(nom.PacketIn{}, &pktInHandler{})
	a.Handle(NewLink{}, &newLinkHandler{})
	a.Handle(lldpTimeout{}, &timeoutHandler{})
	go func() {
		for {
			h.Emit(lldpTimeout{})
			time.Sleep(60 * time.Second)
		}
	}()
}
