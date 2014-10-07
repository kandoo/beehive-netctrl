package controller

import (
	"github.com/soheilhy/beehive-netctrl/nom"
	"github.com/soheilhy/beehive/bh"
)

type pktOutHandler struct{}

func (h *pktOutHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	pkt := msg.Data().(nom.PacketOut)
	d := ctx.Dict(nodeDriversDict)
	k := bh.Key(pkt.Node)
	var nd nodeDrivers
	if err := d.GetGob(k, &nd); err != nil {
		return err
	}

	m := nd.master()
	ctx.SendToBee(pkt, m.BeeID)
	return nil
}

func (h *pktOutHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return bh.MappedCells{
		{nodeDriversDict, bh.Key(msg.Data().(nom.PacketOut).Node)},
	}
}
