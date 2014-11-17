package nom

import (
	"encoding/gob"

	bh "github.com/kandoo/beehive"
)

// Path is a logical sequence of pathlets, where pathlet[i+1] will match the
// output of pathlet[i]. If pathlet[i+1] matches on an incoming port p1,
// pathlet[i] should have a forward action that forwards to a port 2 directly
// connected to p1. Clearly, this rule does not apply to the first and the last
// pathlets in the path.
type Path struct {
	ID       string    // ID needs to be unique only to the subscriber.
	Pathlets []Pathlet // Pathlets in the path.
	Priority uint16    // Priority of this path.
}

func (p Path) Equals(thatp Path) bool {
	if len(p.Pathlets) != len(thatp.Pathlets) {
		return false
	}

	for i := range p.Pathlets {
		if !p.Pathlets[i].Equals(thatp.Pathlets[i]) {
			return false
		}
	}
	return p.ID == thatp.ID
}

// TODO(soheil): add multi-path if there was a real need.

// Pathlet represents a logical connection pathlet in a path, where incoming
// packets matching Match are processing using Actions.
type Pathlet struct {
	Match   Match    // Pathlet's match.
	Exclude []InPort // Exclude packets from these ports in the pathlet.
	Actions []Action // Action that are applied.
}

func (pt Pathlet) Equals(thatpt Pathlet) bool {
	if len(pt.Actions) != len(thatpt.Actions) ||
		len(pt.Exclude) != len(thatpt.Exclude) {

		return false
	}
	if !pt.Match.Equals(thatpt.Match) {
		return false
	}
	for i := range pt.Actions {
		if !pt.Actions[i].Equals(thatpt.Actions[i]) {
			return false
		}
	}
	ports := make(map[InPort]struct{})
	for _, ex := range pt.Exclude {
		ports[ex] = struct{}{}
	}
	for _, ex := range thatpt.Exclude {
		if _, ok := ports[ex]; !ok {
			return false
		}
	}
	return true
}

// AddPath is emitted to install a path in the network.
type AddPath struct {
	Subscriber bh.AppCellKey
	Path       Path
}

// DelPath is emitted to delete a path from the network.
type DelPath struct {
	Path Path
}

// PathAdded is emitted to the subscriber when the path is successfully added.
type PathAdded struct {
	Path Path
}

// PathDeleted is emitted to the subscriber when the path is deleted (because it
// cannot be installed in the network, or because it is explicitly removed).
type PathDeleted struct {
	Path   Path
	Reason PathDelReason
}

// PathDelReason is the reason that a path is deleted.
type PathDelReason int

const (
	// PathDelExplicit means that the path is explicitly deleted using a DelPath.
	PathDelExplicit PathDelReason = iota
	// PathDelInvalid means that the path has contradicting pathlets.
	PathDelInvalid
	// PathDelInfeasible means that the path is valid but cannot be formed due to
	// the current state of the network.
	PathDelInfeasible
)

func init() {
	gob.Register(AddPath{})
	gob.Register(DelPath{})
	gob.Register(Path{})
	gob.Register(PathAdded{})
	gob.Register(PathDeleted{})
	gob.Register(PathDelReason(0))
	gob.Register(Pathlet{})
}
