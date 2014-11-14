package controller

import (
	"encoding/gob"
	"time"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
)

// AddTrigger is a message emitted to install a trigger.
type AddTrigger Trigger

// DelTrigger is a message emitted to remove a trigger for a specific
// subscriber.
type DelTrigger Trigger

// Trigger represents a flow trigger. When a trigger is installed on a node, the
// node will be queried continuiously. Whenever any matching flow goes beyond
// the given bandwidth consumption or lives beyond the given duration, a
// Triggered message will be emitted.
type Trigger struct {
	Subscriber bh.AppCellKey // Triggered messages sent to the Subscriber.
	Node       nom.UID       // The node.
	Match      nom.Match     // The mathing criteria.
	Exact      bool          // Whether Match should exactly match the flow.
	Duration   time.Duration // Minimum live duration to trigger.
	Bandwidth  nom.Bandwidth // Minimum bandwidth consumption to trigger.
}

func (t Trigger) Equals(that Trigger) bool {
	return t.Subscriber == that.Subscriber && t.Node == that.Node &&
		t.Match.Equals(that.Match) && t.Exact == that.Exact &&
		t.Duration == that.Duration && t.Bandwidth == that.Bandwidth
}

// Fired returns whether the trigger is fired according to the stats.
func (t Trigger) Fired(stats nom.FlowStats) bool {
	return t.Bandwidth <= stats.BW() || t.Duration <= stats.Duration
}

// Triggered is a message emitted when a trigger is triggered.
type Triggered struct {
	Node      nom.UID
	Match     nom.Match
	Duration  time.Duration
	Bandwidth nom.Bandwidth
}

type triggerHandler struct{}

func (h triggerHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	add := msg.Data().(AddTrigger)
	var nt nodeTriggers
	ctx.Dict(triggersDict).GetGob(string(add.Node), &nt)
	nt.maybeAddTrigger(Trigger(add))
	return ctx.Dict(triggersDict).PutGob(string(add.Node), &nt)
}

func (h triggerHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return nodeDriversMap(msg.Data().(AddTrigger).Node)
}

func init() {
	gob.Register(AddTrigger{})
	gob.Register(DelTrigger{})
	gob.Register(Trigger{})
}
