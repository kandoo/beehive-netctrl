package controller

import (
	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
)

type addTriggerHandler struct{}

func (h addTriggerHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	add := msg.Data().(nom.AddTrigger)
	dict := ctx.Dict(triggersDict)
	var nt nodeTriggers
	if v, err := dict.Get(string(add.Node)); err == nil {
		nt = v.(nodeTriggers)
	}
	nt.maybeAddTrigger(nom.Trigger(add))
	return dict.Put(string(add.Node), nt)
}

func (h addTriggerHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return nodeDriversMap(msg.Data().(nom.AddTrigger).Node)
}

type delTriggerHandler struct{}

func (h delTriggerHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	del := msg.Data().(nom.DelTrigger)
	var nt nodeTriggers
	dict := ctx.Dict(triggersDict)
	if v, err := dict.Get(string(del.Node)); err == nil {
		nt = v.(nodeTriggers)
	}
	nt.delTrigger(nom.Trigger(del))
	return dict.Put(string(del.Node), nt)
}

func (h delTriggerHandler) Map(msg bh.Msg, ctx bh.RcvContext) bh.MappedCells {
	return nodeDriversMap(msg.Data().(nom.DelTrigger).Node)
}
