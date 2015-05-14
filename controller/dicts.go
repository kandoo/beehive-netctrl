package controller

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"time"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
	"github.com/kandoo/beehive/state"
)

const (
	driversDict  = "ND"
	genDict      = "GD"
	triggersDict = "TD"
	flowsDict    = "FD"
)

type driverInfo struct {
	nom.Driver           // The driver.
	LastSeen   time.Time // The timestamp of the last pong.
	OutPings   int       // Number of pings on the fly.
}

type nodeDrivers struct {
	Node    nom.Node
	Drivers []driverInfo
	Ports   nom.Ports
}

func (nd *nodeDrivers) UID() nom.UID {
	return nd.Node.UID()
}

func (nd *nodeDrivers) GoDecode(b []byte) error {
	return nom.ObjGoDecode(nd, b)
}

func (nd *nodeDrivers) GoEncode() ([]byte, error) {
	return nom.ObjGoEncode(nd)
}

func (nd *nodeDrivers) JSONDecode(b []byte) error {
	return json.Unmarshal(b, nd)
}

func (nd *nodeDrivers) JSONEncode() ([]byte, error) {
	return json.Marshal(nd)
}

func (nd nodeDrivers) driver(d nom.Driver) (driverInfo, bool) {
	for _, e := range nd.Drivers {
		if e.BeeID == d.BeeID {
			return e, true
		}
	}

	return driverInfo{}, false
}

func (nd *nodeDrivers) removeDriver(d nom.Driver) bool {
	for i, e := range nd.Drivers {
		if e.BeeID == d.BeeID {
			nd.Drivers = append(nd.Drivers[:i], nd.Drivers[i+1:]...)
			return true
		}
	}
	return false
}

func (nd *nodeDrivers) master() driverInfo {
	if len(nd.Drivers) == 0 {
		return driverInfo{}
	}

	for _, d := range nd.Drivers {
		if d.Role == nom.DriverRoleMaster {
			return d
		}
	}
	return nd.Drivers[0]
}

func (nd *nodeDrivers) isMaster(d nom.Driver) bool {
	return nd.master().BeeID == d.BeeID
}

func nodeDriversMap(node nom.UID) bh.MappedCells {
	return bh.MappedCells{{driversDict, string(node)}}
}

func sendToMaster(msg interface{}, node nom.UID, ctx bh.RcvContext) error {
	d := ctx.Dict(driversDict)
	var nd nodeDrivers
	if err := d.GetGob(string(node), &nd); err != nil {
		if err == state.ErrNoSuchKey {
			return fmt.Errorf("nom: cannot find node %s", node)
		}
		return err
	}
	ctx.SendToBee(msg, nd.master().BeeID)
	return nil
}

type nodeTriggers struct {
	Node     nom.Node
	Triggers []nom.Trigger
}

func (nt nodeTriggers) hasTrigger(trigger nom.Trigger) bool {
	for _, t := range nt.Triggers {
		if t.Equals(trigger) {
			return true
		}
	}
	return false
}

func (nt *nodeTriggers) maybeAddTrigger(t nom.Trigger) bool {
	if nt.hasTrigger(t) {
		return false
	}
	nt.Triggers = append(nt.Triggers, t)
	return true
}

func (nt *nodeTriggers) delTrigger(t nom.Trigger) {
	panic("todo: implement delTrigger")
}

func newTriggered(t nom.Trigger, d time.Duration,
	bw nom.Bandwidth) nom.Triggered {

	return nom.Triggered{
		Node:      t.Node,
		Match:     t.Match,
		Duration:  d,
		Bandwidth: bw,
	}
}

type flow struct {
	FlowEntry       nom.FlowEntry
	FlowSubscribers []bh.AppCellKey
	Duration        time.Duration
	Packets         uint64
	Bytes           uint64
}

func (f flow) bw() nom.Bandwidth {
	return nom.Bandwidth(f.Bytes / uint64(f.Duration))
}

func (f *flow) updateStats(stats nom.FlowStats) {
	f.Duration = stats.Duration
	f.Bytes = stats.Bytes
	f.Packets = stats.Packets
}

func (f flow) hasFlowSubscriber(sub bh.AppCellKey) bool {
	for _, s := range f.FlowSubscribers {
		if s == sub {
			return true
		}
	}
	return false
}

func (f *flow) maybeAddFlowSubscriber(sub bh.AppCellKey) bool {
	if f.hasFlowSubscriber(sub) {
		return false
	}
	f.FlowSubscribers = append(f.FlowSubscribers, sub)
	return true
}

type nodeFlows struct {
	Node  nom.Node
	Flows []flow
}

func (nf *nodeFlows) flowIndex(flow nom.FlowEntry) int {
	for i := range nf.Flows {
		if flow.Equals(nf.Flows[i].FlowEntry) {
			return i
		}
	}
	return -1
}

func (nf *nodeFlows) maybeAddFlow(add nom.AddFlowEntry) bool {
	i := nf.flowIndex(add.Flow)
	if i < 0 {
		f := flow{
			FlowEntry:       add.Flow,
			FlowSubscribers: []bh.AppCellKey{add.Subscriber},
		}
		nf.Flows = append(nf.Flows, f)
		return true
	}
	return nf.Flows[i].maybeAddFlowSubscriber(add.Subscriber)
}

func init() {
	gob.Register(driverInfo{})
	gob.Register(flow{})
	gob.Register(nodeDrivers{})
	gob.Register(nodeFlows{})
	gob.Register(nodeTriggers{})
}
