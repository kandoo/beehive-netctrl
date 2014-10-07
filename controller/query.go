package controller

import "github.com/soheilhy/beehive/bh"

type queryHandler struct{}

func (h *queryHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	return nil
}

func (h *queryHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return nil
}
