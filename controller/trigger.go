package controller

import (
	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
)

type addTriggerHandler struct{}

func (h addTriggerHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	add := msg.Data().(nom.AddTrigger)
	var nt nodeTriggers
	ctx.Dict(triggersDict).GetGob(string(add.Node), &nt)
	nt.maybeAddTrigger(nom.Trigger(add))
	return ctx.Dict(triggersDict).PutGob(string(add.Node), &nt)
}

func (h addTriggerHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return nodeDriversMap(msg.Data().(nom.AddTrigger).Node)
}
