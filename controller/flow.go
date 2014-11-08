package controller

import (
	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
)

type addFlowHandler struct{}

func (h *addFlowHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	add := msg.Data().(nom.AddFlowEntry)
	return sendToMaster(add, add.Flow.Node, ctx)
}

func (h *addFlowHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return nodeDriversMap(msg.Data().(nom.AddFlowEntry).Flow.Node)
}

type delFlowHandler struct{}

func (h *delFlowHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	return nil
}

func (h *delFlowHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return nil
}
