package controller

import (
	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
)

type pktOutHandler struct{}

func (h *pktOutHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	pkt := msg.Data().(nom.PacketOut)
	d := ctx.Dict(nodeDriversDict)
	k := string(pkt.Node)
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
		{nodeDriversDict, string(msg.Data().(nom.PacketOut).Node)},
	}
}
