package controller

import (
	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
)

type addFlowHandler struct{}

func (h addFlowHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	add := msg.Data().(nom.AddFlowEntry)
	var nf nodeFlows
	if v, err := ctx.Dict(flowsDict).Get(string(add.Flow.Node)); err == nil {
		nf = v.(nodeFlows)
	}
	added := nom.FlowEntryAdded{Flow: add.Flow}
	if nf.maybeAddFlow(add) {
		ctx.Emit(added)
		sendToMaster(add, add.Flow.Node, ctx)
	}
	if !add.Subscriber.IsNil() {
		ctx.SendToCell(added, add.Subscriber.App, add.Subscriber.Cell())
	}
	return ctx.Dict(flowsDict).Put(string(add.Flow.Node), nf)
}

func (h addFlowHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return nodeDriversMap(msg.Data().(nom.AddFlowEntry).Flow.Node)
}

type delFlowHandler struct{}

func (h delFlowHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	return nil
}

func (h delFlowHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return nil
}
