package controller

import "github.com/kandoo/beehive-netctrl/nom"

// Batch represents a collection of flow entries to be added to and removed from
// the network. The controller ensures that these changes are either all applied
// or all failed.
type Batch struct {
	Adds []nom.AddFlowEntry // The flows to be added.
	Dels []nom.DelFlowEntry // The flows to be removed.
}
