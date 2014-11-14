package nom

import (
	"encoding/gob"
	"time"

	bh "github.com/kandoo/beehive"
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
	Node       UID           // The node.
	Match      Match         // The mathing criteria.
	Exact      bool          // Whether Match should exactly match the flow.
	Duration   time.Duration // Minimum live duration to trigger.
	Bandwidth  Bandwidth     // Minimum bandwidth consumption to trigger.
}

func (t Trigger) Equals(that Trigger) bool {
	return t.Subscriber == that.Subscriber && t.Node == that.Node &&
		t.Match.Equals(that.Match) && t.Exact == that.Exact &&
		t.Duration == that.Duration && t.Bandwidth == that.Bandwidth
}

// Fired returns whether the trigger is fired according to the stats.
func (t Trigger) Fired(stats FlowStats) bool {
	return t.Bandwidth <= stats.BW() || t.Duration <= stats.Duration
}

// Triggered is a message emitted when a trigger is triggered.
type Triggered struct {
	Node      UID
	Match     Match
	Duration  time.Duration
	Bandwidth Bandwidth
}

func init() {
	gob.Register(AddTrigger{})
	gob.Register(DelTrigger{})
	gob.Register(Trigger{})
}
