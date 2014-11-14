package controller

import (
	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
)

type addTriggerHandler struct{}

func (h addTriggerHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	add := msg.Data().(nom.AddTrigger)
	var nt nodeTriggers
	dict := ctx.Dict(triggersDict)
	dict.GetGob(string(add.Node), &nt)
	nt.maybeAddTrigger(nom.Trigger(add))
	return dict.PutGob(string(add.Node), &nt)
}

func (h addTriggerHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return nodeDriversMap(msg.Data().(nom.AddTrigger).Node)
}

type delTriggerHandler struct{}

func (h delTriggerHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	del := msg.Data().(nom.DelTrigger)
	var nt nodeTriggers
	dict := ctx.Dict(triggersDict)
	dict.GetGob(string(del.Node), &nt)
	nt.delTrigger(nom.Trigger(del))
	return dict.PutGob(string(del.Node), &nt)
}

func (h delTriggerHandler) Map(msg bh.Msg, ctx bh.RcvContext) bh.MappedCells {
	return nodeDriversMap(msg.Data().(nom.DelTrigger).Node)
}
