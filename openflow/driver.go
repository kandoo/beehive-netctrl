package openflow

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/soheilhy/beehive-netctrl/nom"
	"github.com/soheilhy/beehive-netctrl/openflow/of"
	"github.com/soheilhy/beehive-netctrl/openflow/of10"
	"github.com/soheilhy/beehive-netctrl/openflow/of12"
	"github.com/soheilhy/beehive/bh"
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
		return fmt.Errorf("Received unsupported packet: %v", pkt.Type())
	}
}

func (d *of10Driver) handleMsg(msg bh.Msg, c *ofConn) error {
	ofh, err := d.convToOF(msg)
	if err != nil {
		return fmt.Errorf("of10Driver: Unsupported message %v", msg.Data())
	}

	if err := c.WriteHeader(ofh); err != nil {
		glog.Errorf("ofconn: Cannot write packet: %v", err)
		return err
	}

	return nil
}

func (d *of12Driver) handleMsg(msg bh.Msg, c *ofConn) error {
	ofh, err := d.convToOF(msg)
	if err != nil {
		return fmt.Errorf("of12Driver: Unsupported message %v", msg.Data())
	}

	if err := c.WriteHeader(ofh); err != nil {
		glog.Errorf("ofconn: Cannot write packet: %v", err)
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
			BeeID: c.ctx.BeeID(),
		},
	})
}

func (d *of10Driver) convToOF(msg bh.Msg) (of.Header, error) {
	switch data := msg.Data().(type) {
	case nom.PacketOut:
		ofPort, ok := d.nomPorts[data.InPort]
		if !ok {
			return of.Header{},
				fmt.Errorf("of10Driver: Port %v not found", data.InPort)
		}

		buf := make([]byte, 256)
		out := of10.NewPacketOutWithBuf(buf)
		out.Init()
		out.SetBufferId(uint32(data.BufferID))
		out.SetInPort(ofPort)

		if data.BufferID == 0xFFFFFFFF {
			for _, d := range data.Packet {
				out.AddData(d)
			}
		}

		for _, a := range data.Actions {
			ofa, err := d.convAction(a)
			if err != nil {
				return of.Header{},
					fmt.Errorf("of10Driver: Invalid action %v", err)
			}
			out.AddActions(ofa)
		}

		return out.Header, nil

	default:
		return of.Header{}, fmt.Errorf("of10Driver: Message not supported %v", data)
	}
}

func (d *of12Driver) convToOF(msg bh.Msg) (of.Header, error) {
	return of.Header{}, fmt.Errorf("of12Driver: Message not supported %v",
		msg.Data())
}

func (d *of10Driver) convAction(a nom.Action) (of10.ActionHeader, error) {
	switch action := a.(type) {
	case nom.ActionFlood:
		flood := of10.NewActionOutput()
		flood.SetPort(uint16(of10.PP_FLOOD))
		return flood.ActionHeader, nil

	default:
		return of10.ActionHeader{},
			fmt.Errorf("of10Driver: Action not support %v", action)
	}
}
