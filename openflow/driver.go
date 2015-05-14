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
	case of10.IsStatsReply(pkt10):
		return d.handleStatsReply(of10.NewStatsReplyWithBuf(pkt10.Buf), c)
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
	case of12.IsStatsReply(pkt12):
		return d.handleStatsReply(of12.NewStatsReplyWithBuf(pkt12.Buf), c)
	case of12.IsRoleReply(pkt12):
		return d.handleRoleReply(of12.NewRoleReplyWithBuf(pkt12.Buf), c)
	default:
		return fmt.Errorf("received unsupported packet: %v", pkt.Type())
	}
}

func (d *of10Driver) handleMsg(msg bh.Msg, c *ofConn) error {
	ofh, err := d.convToOF(msg, c)
	if err != nil {
		return err
	}

	// Means we should ignore the OpenFlow header.
	if ofh.Size() == 0 {
		return nil
	}

	if err := c.WriteHeader(ofh); err != nil {
		glog.Errorf("ofconn: cannot write packet: %v", err)
		return err
	}

	return nil
}

func (d *of12Driver) handleMsg(msg bh.Msg, c *ofConn) error {
	ofh, err := d.convToOF(msg, c)
	if err != nil {
		return err
	}

	// Means we should ignore the OpenFlow header.
	if ofh.Size() == 0 {
		return nil
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

func (d *of10Driver) convToOF(msg bh.Msg, c *ofConn) (of.Header, error) {
	switch data := msg.Data().(type) {
	case nom.Ping:
		c.ctx.ReplyTo(msg, nom.Pong{})
		return of.Header{}, nil

	case nom.ChangeDriverRole:
		if data.Role == nom.DriverRoleSlave {
			return of.Header{}, errors.New("of10Driver: role slave is not supported")
		}

		return of.Header{}, nil

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
		match, err := d.ofMatch(data.Flow.Match)
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
		match, err := d.ofMatch(data.Match)
		if err != nil {
			return of.Header{}, fmt.Errorf("of10Driver: invalid match %v", err)
		}
		mod.SetMatch(match)
		return mod.Header, nil

	case nom.FlowStatsQuery:
		match, err := d.ofMatch(data.Match)
		if err != nil {
			return of.Header{}, fmt.Errorf("of10Driver: invalid match %v", err)
		}
		query := of10.NewFlowStatsRequest()
		query.SetMatch(match)
		query.SetTableId(0xFF)
		query.SetOutPort(uint16(of10.PP_NONE))
		return query.Header, nil

	default:
		return of.Header{}, fmt.Errorf("of10Driver: unsupported message %#v", data)
	}
}

func (d *of12Driver) convToOF(msg bh.Msg, c *ofConn) (of.Header, error) {
	switch data := msg.Data().(type) {
	case nom.Ping:
		c.ctx.ReplyTo(msg, nom.Pong{})
		return of.Header{}, nil

	case nom.ChangeDriverRole:
		rr := of12.NewRoleRequest()
		rr.SetGenerationId(data.Generation)
		switch data.Role {
		case nom.DriverRoleDefault:
			rr.SetRole(uint32(of12.ROLE_EQUAL))
		case nom.DriverRoleMaster:
			rr.SetRole(uint32(of12.ROLE_MASTER))
		case nom.DriverRoleSlave:
			rr.SetRole(uint32(of12.ROLE_SLAVE))
		}
		return rr.Header, nil

	case nom.PacketOut:
		buf := make([]byte, 256)
		out := of12.NewPacketOutWithBuf(buf)
		out.Init()
		out.SetBufferId(uint32(data.BufferID))

		ofPort, ok := d.nomPorts[data.InPort]
		if ok {
			out.SetInPort(ofPort)
		}

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
		mod := of12.NewFlowMod()
		mod.SetCommand(uint8(of12.PFC_ADD))
		mod.SetPriority(uint16(data.Flow.Priority))
		mod.SetIdleTimeout(uint16(data.Flow.IdleTimeout))
		mod.SetHardTimeout(uint16(data.Flow.HardTimeout))
		mod.SetBufferId(^uint32(0))
		match, err := d.ofMatch(data.Flow.Match)
		if err != nil {
			return of.Header{}, fmt.Errorf("of12Driver: invalid match %v", err)
		}
		mod.SetMatch(match)

		inst := of12.NewApplyActions()
		for _, a := range data.Flow.Actions {
			ofa, err := d.convAction(a)
			if err != nil {
				return of.Header{},
					fmt.Errorf("of12Driver: invalid action %v", err)
			}
			inst.AddActions(ofa)
		}
		mod.AddInstructions(inst.Instruction)

		return mod.Header, nil

	case nom.DelFlowEntry:
		mod := of12.NewFlowMod()
		if data.Exact {
			mod.SetCommand(uint8(of10.PFC_DELETE))
		} else {
			mod.SetCommand(uint8(of10.PFC_DELETE_STRICT))
		}
		match, err := d.ofMatch(data.Match)
		if err != nil {
			return of.Header{}, fmt.Errorf("of10Driver: invalid match %v", err)
		}
		mod.SetMatch(match)
		return mod.Header, nil

	case nom.FlowStatsQuery:
		match, err := d.ofMatch(data.Match)
		if err != nil {
			return of.Header{}, fmt.Errorf("of10Driver: invalid match %v", err)
		}
		query := of12.NewFlowStatsRequest()
		query.SetTableId(0xFF)
		query.SetOutPort(uint32(of12.PP_ANY))
		query.SetOutGroup(uint32(of12.PG_ANY))
		query.SetMatch(match)
		return query.Header, nil

	default:
		return of.Header{}, fmt.Errorf("of12Driver: unsupported message %#v", data)
	}
}

func (d *of10Driver) convAction(a nom.Action) (of10.Action, error) {
	switch action := a.(type) {
	case nom.ActionFlood:
		flood := of10.NewActionOutput()
		flood.SetPort(uint16(of10.PP_FLOOD))
		return flood.Action, nil

	case nom.ActionForward:
		if len(action.Ports) != 1 {
			return of10.Action{},
				errors.New("of10Driver: can forward to only one port")
		}
		p, ok := d.nomPorts[action.Ports[0]]
		if !ok {
			return of10.Action{},
				fmt.Errorf("of10Driver: port %v no found", action.Ports[0])
		}
		out := of10.NewActionOutput()
		out.SetPort(uint16(p))
		return out.Action, nil

	default:
		return of10.Action{},
			fmt.Errorf("of10Driver: action not supported %v", action)
	}
}

func (d *of10Driver) nomMatch(m of10.Match) (nom.Match, error) {
	nm := nom.Match{}
	wc := of10.FlowWildcards(m.Wildcards())
	if wc&of10.PFW_IN_PORT == 0 {
		ofp := m.InPort()
		nomp, ok := d.ofPorts[ofp]
		if !ok {
			return nom.Match{}, fmt.Errorf("of10Driver: cannot find port %v", ofp)
		}
		nm.AddField(nom.InPort(nomp.UID()))
	}
	if wc&of10.PFW_DL_SRC == 0 {
		nm.AddField(nom.EthSrc{
			Addr: m.DlSrc(),
			Mask: nom.MaskNoneMAC,
		})
	}
	if wc&of10.PFW_DL_DST == 0 {
		nm.AddField(nom.EthDst{
			Addr: m.DlDst(),
			Mask: nom.MaskNoneMAC,
		})
	}
	if wc&of10.PFW_DL_TYPE == 0 {
		nm.AddField(nom.EthType(m.DlType()))
	}
	if wc&of10.PFW_NW_SRC_MASK != of10.PFW_NW_SRC_ALL && m.NwSrc() != 0 {
		mask := uint(wc&of10.PFW_NW_SRC_MASK) >> uint(of10.PFW_NW_SRC_SHIFT)
		nm.AddField(nom.IPv4Src(nom.CIDRToMaskedIPv4(m.NwSrc(), mask)))
	}
	if wc&of10.PFW_NW_DST_MASK != of10.PFW_NW_DST_ALL && m.NwDst() != 0 {
		mask := uint(wc&of10.PFW_NW_DST_MASK) >> uint(of10.PFW_NW_DST_SHIFT)
		nm.AddField(nom.IPv4Dst(nom.CIDRToMaskedIPv4(m.NwDst(), mask)))
	}
	if wc&of10.PFW_TP_SRC == 0 {
		nm.AddField(nom.TransportPortSrc(m.TpSrc()))
	}
	if wc&of10.PFW_TP_DST == 0 {
		nm.AddField(nom.TransportPortDst(m.TpDst()))
	}
	return nm, nil
}

func (d *of10Driver) ofMatch(m nom.Match) (of10.Match, error) {
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
			if f.Mask != nom.MaskNoneMAC {
				return of10.Match{},
					fmt.Errorf("of10Driver: masked ethernet address is not supported")
			}
			ofm.SetDlDst([6]byte(f.Addr))
			w &= ^of10.PFW_DL_DST

		case nom.EthSrc:
			if f.Mask != nom.MaskNoneMAC {
				return of10.Match{},
					fmt.Errorf("of10Driver: masked ethernet address is not supported")
			}
			ofm.SetDlSrc([6]byte(f.Addr))
			w &= ^of10.PFW_DL_SRC

		case nom.EthType:
			ofm.SetDlType(uint16(f))
			w &= ^of10.PFW_DL_TYPE

		case nom.IPv4Src:
			w &= ^of10.PFW_NW_SRC_MASK
			ofm.SetNwSrc(f.Addr.Uint())
			mask := f.Mask.PopCount()
			if mask == 32 {
				w |= of10.PFW_NW_SRC_MASK
			} else {
				w |= of10.FlowWildcards(mask << uint(of10.PFW_NW_SRC_SHIFT))
			}

		case nom.IPv4Dst:
			w &= ^of10.PFW_NW_DST_MASK
			ofm.SetNwDst(f.Addr.Uint())
			mask := f.Mask.PopCount()
			if mask == 32 {
				w |= of10.PFW_NW_DST_MASK
			} else {
				w |= of10.FlowWildcards(mask << uint(of10.PFW_NW_DST_SHIFT))
			}

		case nom.TransportPortSrc:
			ofm.SetTpSrc(uint16(f))
			w &= ^of10.PFW_TP_SRC

		case nom.TransportPortDst:
			ofm.SetTpDst(uint16(f))
			w &= ^of10.PFW_TP_DST

		case nom.IPv6Src, nom.IPv6Dst:
			return of10.Match{}, fmt.Errorf("of10Driver: IPv6 not supported")
		}
	}
	ofm.SetWildcards(uint32(w))
	return ofm, nil
}

func (d *of12Driver) convAction(a nom.Action) (of12.Action, error) {
	switch action := a.(type) {
	case nom.ActionFlood:
		flood := of12.NewActionOutput()
		flood.SetPort(uint32(of12.PP_FLOOD))
		return flood.Action, nil

	case nom.ActionForward:
		if len(action.Ports) != 1 {
			return of12.Action{},
				errors.New("of12Driver: can forward to only one port")
		}
		p, ok := d.nomPorts[action.Ports[0]]
		if !ok {
			return of12.Action{},
				fmt.Errorf("of12Driver: port %v no found", action.Ports[0])
		}
		out := of12.NewActionOutput()
		out.SetPort(uint32(p))
		return out.Action, nil

	default:
		return of12.Action{},
			fmt.Errorf("of12Driver: action not supported %v", action)
	}
}

func (d *of12Driver) ofMatch(m nom.Match) (of12.Match, error) {
	ofm := of12.NewOXMatch()
	for _, f := range m.Fields {
		switch f := f.(type) {
		case nom.InPort:
			p, ok := d.nomPorts[nom.UID(f)]
			if !ok {
				return of12.Match{}, fmt.Errorf("of12Driver: nom port not found %v", f)
			}
			off := of12.NewOxmInPort()
			off.SetInPort(p)
			ofm.AddFields(off.OxmField)

		case nom.EthDst:
			if f.Mask == nom.MaskNoneMAC {
				off := of12.NewOxmEthDst()
				off.SetMacAddr([6]byte(f.Addr))
				ofm.AddFields(off.OxmField)
			} else {
				off := of12.NewOxmEthDstMasked()
				off.SetMacAddr([6]byte(f.Addr))
				off.SetMask([6]byte(f.Mask))
				ofm.AddFields(off.OxmField)
			}

		case nom.EthSrc:
			if f.Mask == nom.MaskNoneMAC {
				off := of12.NewOxmEthSrc()
				off.SetMacAddr([6]byte(f.Addr))
				ofm.AddFields(off.OxmField)
			} else {
				off := of12.NewOxmEthSrcMasked()
				off.SetMacAddr([6]byte(f.Addr))
				off.SetMask([6]byte(f.Mask))
				ofm.AddFields(off.OxmField)
			}

		case nom.EthType:
			off := of12.NewOxmEthType()
			off.SetType(uint16(f))
			ofm.AddFields(off.OxmField)

		case nom.IPProto:
			off := of12.NewOxmIpProto()
			off.SetProto(uint8(f))
			ofm.AddFields(off.OxmField)

		case nom.IPv4Src:
			if f.Mask == nom.MaskNoneIPV4 {
				off := of12.NewOxmIpV4Src()
				off.SetAddr(f.Addr)
				ofm.AddFields(off.OxmField)
			} else {
				off := of12.NewOxmIpV4SrcMasked()
				off.SetAddr(f.Addr)
				off.SetMask(f.Mask)
				ofm.AddFields(off.OxmField)
			}

		case nom.IPv4Dst:
			if f.Mask == nom.MaskNoneIPV4 {
				off := of12.NewOxmIpV4Dst()
				off.SetAddr(f.Addr)
				ofm.AddFields(off.OxmField)
			} else {
				off := of12.NewOxmIpV4DstMasked()
				off.SetAddr(f.Addr)
				off.SetMask(f.Mask)
				ofm.AddFields(off.OxmField)
			}

		case nom.IPv6Src:
			if f.Mask == nom.MaskNoneIPV6 {
				off := of12.NewOxmIpV6Src()
				off.SetAddr(f.Addr)
				ofm.AddFields(off.OxmField)
			} else {
				off := of12.NewOxmIpV6SrcMasked()
				off.SetAddr(f.Addr)
				off.SetMask(f.Mask)
				ofm.AddFields(off.OxmField)
			}

		case nom.IPv6Dst:
			if f.Mask == nom.MaskNoneIPV6 {
				off := of12.NewOxmIpV6Dst()
				off.SetAddr(f.Addr)
				ofm.AddFields(off.OxmField)
			} else {
				off := of12.NewOxmIpV6DstMasked()
				off.SetAddr(f.Addr)
				off.SetMask(f.Mask)
				ofm.AddFields(off.OxmField)
			}

		case nom.TransportPortSrc:
			off := of12.NewOxmTcpSrc()
			off.SetPort(uint16(f))
			ofm.AddFields(off.OxmField)

		case nom.TransportPortDst:
			off := of12.NewOxmTcpDst()
			off.SetPort(uint16(f))
			ofm.AddFields(off.OxmField)

		default:
			return of12.Match{}, fmt.Errorf("of12Driver: %#v is not supported", f)
		}
	}
	return ofm.Match, nil
}

func (d *of12Driver) nomMatch(m of12.Match) (nom.Match, error) {
	nm := nom.Match{}

	if !of12.IsOXMatch(m) {
		return nm, fmt.Errorf("of12Driver: std math is not supported")
	}

	xm, err := of12.ToOXMatch(m)
	if err != nil {
		return nm, err
	}

	for _, f := range xm.Fields() {
		switch f.OxmField() {
		case uint8(of12.PXMT_IN_PORT):
			xf, err := of12.ToOxmInPort(f)
			if err != nil {
				return nom.Match{}, err
			}

			np, ok := d.ofPorts[xf.InPort()]
			if !ok {
				return nom.Match{}, fmt.Errorf("of12Driver: cannot find port %v",
					xf.InPort())
			}
			nm.AddField(nom.InPort(np.UID()))

		case uint8(of12.PXMT_ETH_TYPE):
			xf, err := of12.ToOxmEthType(f)
			if err != nil {
				return nom.Match{}, err
			}
			nm.AddField(nom.EthType(xf.Type()))

		case uint8(of12.PXMT_ETH_SRC):
			xf, err := of12.ToOxmEthSrc(f)
			if err != nil {
				return nom.Match{}, err
			}

			nf := nom.EthSrc{}
			nf.Addr = xf.MacAddr()
			nf.Mask = nom.MaskNoneMAC
			nm.AddField(nf)

		case uint8(of12.PXMT_ETH_SRC_MASKED):
			xf, err := of12.ToOxmEthSrcMasked(f)
			if err != nil {
				return nom.Match{}, err
			}

			nf := nom.EthSrc{}
			nf.Addr = xf.MacAddr()
			nf.Mask = xf.Mask()
			nm.AddField(nf)

		case uint8(of12.PXMT_ETH_DST):
			xf, err := of12.ToOxmEthDst(f)
			if err != nil {
				return nom.Match{}, err
			}

			nf := nom.EthDst{}
			nf.Addr = xf.MacAddr()
			nf.Mask = nom.MaskNoneMAC
			nm.AddField(nf)

		case uint8(of12.PXMT_ETH_DST_MASKED):
			xf, err := of12.ToOxmEthDstMasked(f)
			if err != nil {
				return nom.Match{}, err
			}

			nf := nom.EthDst{}
			nf.Addr = xf.MacAddr()
			nf.Mask = xf.Mask()
			nm.AddField(nf)

		case uint8(of12.PXMT_IP_PROTO):
			xf, err := of12.ToOxmIpProto(f)
			if err != nil {
				return nom.Match{}, err
			}

			nm.AddField(nom.IPProto(xf.Proto()))

		case uint8(of12.PXMT_IPV4_SRC):
			xf, err := of12.ToOxmIpV4Src(f)
			if err != nil {
				return nom.Match{}, err
			}

			nf := nom.IPv4Src{}
			nf.Addr = nom.IPv4Addr(xf.Addr())
			nf.Mask = nom.MaskNoneIPV4
			nm.AddField(nf)

		case uint8(of12.PXMT_IPV4_SRC_MASKED):
			xf, err := of12.ToOxmIpV4SrcMasked(f)
			if err != nil {
				return nom.Match{}, err
			}

			nf := nom.IPv4Src{}
			nf.Addr = nom.IPv4Addr(xf.Addr())
			nf.Mask = nom.IPv4Addr(xf.Mask())
			nm.AddField(nf)

		case uint8(of12.PXMT_IPV4_DST):
			xf, err := of12.ToOxmIpV4Dst(f)
			if err != nil {
				return nom.Match{}, err
			}

			nf := nom.IPv4Dst{}
			nf.Addr = nom.IPv4Addr(xf.Addr())
			nf.Mask = nom.MaskNoneIPV4
			nm.AddField(nf)

		case uint8(of12.PXMT_IPV4_DST_MASKED):
			xf, err := of12.ToOxmIpV4DstMasked(f)
			if err != nil {
				return nom.Match{}, err
			}

			nf := nom.IPv4Dst{}
			nf.Addr = nom.IPv4Addr(xf.Addr())
			nf.Mask = nom.IPv4Addr(xf.Mask())
			nm.AddField(nf)

		case uint8(of12.PXMT_IPV6_SRC):
			xf, err := of12.ToOxmIpV6Src(f)
			if err != nil {
				return nom.Match{}, err
			}

			nf := nom.IPv6Src{}
			nf.Addr = nom.IPv6Addr(xf.Addr())
			nf.Mask = nom.MaskNoneIPV6
			nm.AddField(nf)

		case uint8(of12.PXMT_IPV6_SRC_MASKED):
			xf, err := of12.ToOxmIpV6SrcMasked(f)
			if err != nil {
				return nom.Match{}, err
			}

			nf := nom.IPv6Src{}
			nf.Addr = nom.IPv6Addr(xf.Addr())
			nf.Mask = nom.IPv6Addr(xf.Mask())
			nm.AddField(nf)

		case uint8(of12.PXMT_IPV6_DST):
			xf, err := of12.ToOxmIpV6Dst(f)
			if err != nil {
				return nom.Match{}, err
			}

			nf := nom.IPv6Dst{}
			nf.Addr = nom.IPv6Addr(xf.Addr())
			nf.Mask = nom.MaskNoneIPV6
			nm.AddField(nf)

		case uint8(of12.PXMT_IPV6_DST_MASKED):
			xf, err := of12.ToOxmIpV6DstMasked(f)
			if err != nil {
				return nom.Match{}, err
			}

			nf := nom.IPv6Dst{}
			nf.Addr = nom.IPv6Addr(xf.Addr())
			nf.Mask = nom.IPv6Addr(xf.Mask())
			nm.AddField(nf)

		case uint8(of12.PXMT_TCP_SRC):
			xf, err := of12.ToOxmTcpSrc(f)
			if err != nil {
				return nom.Match{}, err
			}

			nm.AddField(nom.TransportPortSrc(xf.Port()))

		case uint8(of12.PXMT_TCP_DST):
			xf, err := of12.ToOxmTcpDst(f)
			if err != nil {
				return nom.Match{}, err
			}

			nm.AddField(nom.TransportPortDst(xf.Port()))
		}
	}

	return nm, nil
}
