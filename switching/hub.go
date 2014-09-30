package switching

import (
	"github.com/soheilhy/beehive-netctrl/nom"
	"github.com/soheilhy/beehive/bh"
)

type Hub struct{}

func (h *Hub) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	in := msg.Data().(nom.PacketIn)
	out := nom.PacketOut{
		Node:     in.Node,
		InPort:   in.InPort,
		BufferID: in.BufferID,
		Packet:   in.Packet,
		Actions:  []nom.Action{nom.ActionFlood{}},
	}
	ctx.ReplyTo(msg, out)
	return nil
}

func (h *Hub) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return bh.MappedCells{{"N", bh.Key(msg.Data().(nom.PacketIn).Node)}}
}
