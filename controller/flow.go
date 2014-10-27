package controller

import bh "github.com/kandoo/beehive"

type addFlowHandler struct{}

func (h *addFlowHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	return nil
}

func (h *addFlowHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return nil
}

type delFlowHandler struct{}

func (h *delFlowHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	return nil
}

func (h *delFlowHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return nil
}
