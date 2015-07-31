package discovery

import (
	"fmt"
	"time"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/net/ethernet"
	"github.com/kandoo/beehive-netctrl/nom"
	"github.com/kandoo/beehive/Godeps/_workspace/src/github.com/golang/glog"
)

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

func (np *nodePortsAndLinks) linkFrom(from nom.UID) (nom.Link, bool) {
	for _, l := range np.L {
		if l.From == from {
			return l, true
		}
	}
	return nom.Link{}, false
}

func (np *nodePortsAndLinks) hasLinkFrom(from nom.UID) bool {
	_, ok := np.linkFrom(from)
	return ok
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
	if v, err := d.Get(k); err != nil {
		glog.Warningf("%v rejoins", n)
	} else {
		np = v.(nodePortsAndLinks)
	}
	np.N = n
	// TODO(soheil): Add a flow entry to forward lldp packets to the controller.
	return d.Put(k, np)
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
	v, err := d.Get(k)
	if err != nil {
		glog.Warningf("%v added before its node", p)
		ctx.Snooze(1 * time.Second)
		return nil
	}

	np := v.(nodePortsAndLinks)
	if np.hasPort(p) {
		glog.Warningf("%v readded")
		np.removePort(p)
	}

	sendLLDPPacket(np.N, p, ctx)

	np.P = append(np.P, p)
	return d.Put(k, np)
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
	d.ForEach(func(k string, v interface{}) bool {
		np := v.(nodePortsAndLinks)
		for _, p := range np.P {
			sendLLDPPacket(np.N, p, ctx)
		}
		return true
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
	if _, err := d.Get(k); err != nil {
		return fmt.Errorf("Node %v not found", pin.Node)
	}

	l := nom.Link{
		ID:    nom.LinkID(port.UID()),
		From:  pin.InPort,
		To:    port.UID(),
		State: nom.LinkStateUp,
	}
	ctx.Emit(NewLink(l))

	l = nom.Link{
		ID:    nom.LinkID(pin.InPort),
		From:  port.UID(),
		To:    pin.InPort,
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
	v, err := d.Get(k)
	if err != nil {
		return err
	}
	np := v.(nodePortsAndLinks)

	if oldl, ok := np.linkFrom(l.From); ok {
		if oldl.UID() == l.UID() {
			return nil
		}
		np.removeLink(oldl)
		ctx.Emit(nom.LinkDeleted(oldl))
	}

	glog.V(2).Infof("Link detected %v", l)
	ctx.Emit(nom.LinkAdded(l))
	np.L = append(np.L, l)
	return d.Put(k, np)
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
