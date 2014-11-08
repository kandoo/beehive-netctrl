package openflow

import (
	"errors"
	"fmt"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
	"github.com/kandoo/beehive-netctrl/openflow/of"
	"github.com/kandoo/beehive-netctrl/openflow/of10"
	"github.com/kandoo/beehive-netctrl/openflow/of12"
	"github.com/kandoo/beehive/Godeps/_workspace/src/github.com/golang/glog"
)

type ofDriver interface {
	handshake(conn *ofConn) error
	handlePkt(pkt of.Header, conn *ofConn) error
	handleMsg(msg bh.Msg, conn *ofConn) error
	handleConnClose(conn *ofConn)
}

type of10Driver struct {
	ofPorts  map[uint16]*nom.Port
	nomPorts map[nom.UID]uint16
}

type of12Driver struct {
	ofPorts  map[uint32]*nom.Port
	nomPorts map[nom.UID]uint32
}

func (d *of10Driver) handlePkt(pkt of.Header, c *ofConn) error {
	pkt10, err := of10.ToHeader10(pkt)
	if err != nil {
		return err
	}

	switch {
	case of10.IsEchoRequest(pkt10):
		return d.handleEchoRequest(of10.NewEchoRequestWithBuf(pkt10.Buf), c)
	case of10.IsFeaturesReply(pkt10):
		return d.handleFeaturesReply(of10.NewFeaturesReplyWithBuf(pkt10.Buf), c)
	case of10.IsPacketIn(pkt10):
		return d.handlePacketIn(of10.NewPacketInWithBuf(pkt10.Buf), c)
	case of10.IsErrorMsg(pkt10):
		return d.handleErrorMsg(of10.NewErrorMsgWithBuf(pkt10.Buf), c)
	default:
		return fmt.Errorf("Received unsupported packet: %v", pkt.Type())
	}
}

func (d *of12Driver) handlePkt(pkt of.Header, c *ofConn) error {
	pkt12, err := of12.ToHeader12(pkt)
	if err != nil {
		return err
	}

	switch {
	case of12.IsEchoRequest(pkt12):
		return d.handleEchoRequest(of12.NewEchoRequestWithBuf(pkt12.Buf), c)
	case of12.IsFeaturesReply(pkt12):
		return d.handleFeaturesReply(of12.NewFeaturesReplyWithBuf(pkt12.Buf), c)
	case of12.IsPacketIn(pkt12):
		return d.handlePacketIn(of12.NewPacketInWithBuf(pkt12.Buf), c)
	case of12.IsErrorMsg(pkt12):
		return d.handleErrorMsg(of12.NewErrorMsgWithBuf(pkt12.Buf), c)
	default:
		return fmt.Errorf("received unsupported packet: %v", pkt.Type())
	}
}

func (d *of10Driver) handleMsg(msg bh.Msg, c *ofConn) error {
	ofh, err := d.convToOF(msg)
	if err != nil {
		return err
	}

	if err := c.WriteHeader(ofh); err != nil {
		glog.Errorf("ofconn: cannot write packet: %v", err)
		return err
	}

	return nil
}

func (d *of12Driver) handleMsg(msg bh.Msg, c *ofConn) error {
	ofh, err := d.convToOF(msg)
	if err != nil {
		return err
	}

	if err := c.WriteHeader(ofh); err != nil {
		glog.Errorf("ofconn: cannot write packet: %v", err)
		return err
	}

	return nil
}

func (d *of10Driver) handleConnClose(c *ofConn) {
	emitNodeDisconnected(c)
}

func (d *of12Driver) handleConnClose(c *ofConn) {
	emitNodeDisconnected(c)
}

func emitNodeDisconnected(c *ofConn) {
	c.ctx.Emit(nom.NodeDisconnected{
		Node: c.node,
		Driver: nom.Driver{
			BeeID: c.ctx.ID(),
		},
	})
}

func (d *of10Driver) convToOF(msg bh.Msg) (of.Header, error) {
	switch data := msg.Data().(type) {
	case nom.PacketOut:
		buf := make([]byte, 256)
		out := of10.NewPacketOutWithBuf(buf)
		out.Init()
		out.SetBufferId(uint32(data.BufferID))

		ofPort, ok := d.nomPorts[data.InPort]
		if ok {
			out.SetInPort(ofPort)
		}

		// FIXME(soheil): when actions are added after data, the packet becomes
		// corrupted.
		for _, a := range data.Actions {
			ofa, err := d.convAction(a)
			if err != nil {
				return of.Header{},
					fmt.Errorf("of10Driver: invalid action %v", err)
			}
			out.AddActions(ofa)
		}

		if data.BufferID == 0xFFFFFFFF {
			for _, d := range data.Packet {
				out.AddData(d)
			}
		}

		return out.Header, nil

	case nom.AddFlowEntry:
		mod := of10.NewFlowMod()
		mod.SetCommand(uint16(of10.PFC_ADD))
		mod.SetPriority(uint16(data.Flow.Priority))
		mod.SetIdleTimeout(uint16(data.Flow.IdleTimeout))
		mod.SetHardTimeout(uint16(data.Flow.HardTimeout))
		mod.SetBufferId(0xFFFFFFFF)
		match, err := d.convMatch(data.Flow.Match)
		if err != nil {
			return of.Header{}, fmt.Errorf("of10Driver: invalid match %v", err)
		}
		mod.SetMatch(match)
		for _, a := range data.Flow.Actions {
			ofa, err := d.convAction(a)
			if err != nil {
				return of.Header{},
					fmt.Errorf("of10Driver: invalid action %v", err)
			}
			mod.AddActions(ofa)
		}
		return mod.Header, nil

	case nom.DelFlowEntry:
		mod := of10.NewFlowMod()
		if data.Exact {
			mod.SetCommand(uint16(of10.PFC_DELETE))
		} else {
			mod.SetCommand(uint16(of10.PFC_DELETE_STRICT))
		}
		match, err := d.convMatch(data.Match)
		if err != nil {
			return of.Header{}, fmt.Errorf("of10Driver: invalid match %v", err)
		}
		mod.SetMatch(match)
		return mod.Header, nil

	default:
		return of.Header{}, fmt.Errorf("of10Driver: unsupported message %+v", data)
	}
}

func (d *of12Driver) convToOF(msg bh.Msg) (of.Header, error) {
	return of.Header{}, fmt.Errorf("of12Driver: message not supported %+v",
		msg.Data())
}

func (d *of10Driver) convAction(a nom.Action) (of10.ActionHeader, error) {
	switch action := a.(type) {
	case nom.ActionFlood:
		flood := of10.NewActionOutput()
		flood.SetPort(uint16(of10.PP_FLOOD))
		return flood.ActionHeader, nil

	case nom.ActionForward:
		if len(action.Ports) != 1 {
			return of10.ActionHeader{},
				errors.New("of10Driver: can forward to only one port")
		}
		p, ok := d.nomPorts[action.Ports[0]]
		if !ok {
			return of10.ActionHeader{},
				fmt.Errorf("of10Driver: port %v no found", action.Ports[0])
		}
		out := of10.NewActionOutput()
		out.SetPort(uint16(p))
		return out.ActionHeader, nil

	default:
		return of10.ActionHeader{},
			fmt.Errorf("of10Driver: action not supported %v", action)
	}
}

func (d *of10Driver) convMatch(m nom.Match) (of10.Match, error) {
	ofm := of10.NewMatch()
	w := of10.PFW_ALL
	for _, f := range m.Fields {
		switch f := f.(type) {
		case nom.InPort:
			p, ok := d.nomPorts[nom.UID(f)]
			if !ok {
				return of10.Match{}, fmt.Errorf("of10Driver: nom port not found %v", f)
			}
			ofm.SetInPort(p)
			w &= ^of10.PFW_IN_PORT

		case nom.EthDst:
			if f.Mask != [6]byte{} {
				return of10.Match{},
					fmt.Errorf("of10Driver: masked ethernet address is not supported")
			}
			ofm.SetDlDst([6]byte(f.Addr))
			w &= ^of10.PFW_DL_DST

		case nom.EthSrc:
			if f.Mask != [6]byte{} {
				return of10.Match{},
					fmt.Errorf("of10Driver: masked ethernet address is not supported")
			}
			ofm.SetDlSrc([6]byte(f.Addr))
			w &= ^of10.PFW_DL_SRC

		case nom.EthType:
			ofm.SetDlType(uint16(f))
			w &= ^of10.PFW_DL_TYPE

		case nom.IPV4Src:
			ofm.SetNwSrc(f.Addr.Uint32())
			mask := f.Mask.Uint32()
			w &= ^of10.PFW_NW_SRC_ALL
			for i := uint(0); i < 32; i++ {
				if mask&(1<<i) != 0 {
					w |= of10.FlowWildcards(i << uint(of10.PFW_NW_SRC_SHIFT))
					break
				}
			}

		case nom.IPV4Dst:
			ofm.SetNwDst(f.Addr.Uint32())
			mask := f.Mask.Uint32()
			w &= ^of10.PFW_NW_DST_ALL
			for i := uint(0); i < 32; i++ {
				if mask&(1<<i) != 0 {
					w |= of10.FlowWildcards(i << uint(of10.PFW_NW_DST_SHIFT))
					break
				}
			}

		case nom.TransportPortSrc:
			ofm.SetTpSrc(uint16(f))
			w &= ^of10.PFW_TP_SRC

		case nom.TransportPortDst:
			ofm.SetTpDst(uint16(f))
			w &= ^of10.PFW_TP_DST

		case nom.IPV6Src, nom.IPV6Dst:
			return of10.Match{}, fmt.Errorf("of10Driver: IPv6 not supported")
		}
	}
	ofm.SetWildcards(uint32(w))
	return ofm, nil
}
