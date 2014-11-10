package nom

import (
	"encoding/gob"
	"time"
)

// NodeQuery queries the information of a node.
type NodeQuery struct {
	Node UID
}

// NodeQueryResult is the result for NodeQuery.
type NodeQueryResult struct {
	Err  error
	Node Node
}

// PortQuery queries the information of a port.
type PortQuery struct {
	Port UID
}

// PortQueryResult is the result for a PortQuery.
type PortQueryResult struct {
	Err  error
	Port Port
}

// FlowStatQuery queries the flows that would match the query. If Exact is
// false, it removes all flow entries that are subsumed by the given match.
type FlowStatQuery struct {
	Match Match
	Exact bool
}

// FlowStatQueryResult is the result for a FlowStatQuery
type FlowStatQueryResult struct {
	Flow     FlowEntry
	Duration time.Duration
	PktCount uint64
	Bytes    uint64
}

func init() {
	gob.Register(FlowStatQuery{})
	gob.Register(FlowStatQueryResult{})
	gob.Register(NodeQuery{})
	gob.Register(NodeQueryResult{})
	gob.Register(PortQuery{})
	gob.Register(PortQueryResult{})
}
