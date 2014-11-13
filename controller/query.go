package controller

import (
	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
)

type queryHandler struct{}

func (h queryHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	query := msg.Data().(nom.FlowStatsQuery)
	return sendToMaster(query, query.Node, ctx)
}

func (h queryHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return nodeDriversMap(msg.Data().(nom.FlowStatsQuery).Node)
}
