package controller

import (
	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
)

type pktOutHandler struct{}

func (h *pktOutHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	pkt := msg.Data().(nom.PacketOut)
	return sendToMaster(pkt, pkt.Node, ctx)
}

func (h *pktOutHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return nodeDriversMap(msg.Data().(nom.PacketOut).Node)
}
