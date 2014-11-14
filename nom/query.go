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

// FlowStatsQuery queries the flows that would match the query. If Exact is
// false, it removes all flow entries that are subsumed by the given match.
type FlowStatsQuery struct {
	Node  UID
	Match Match
}

// FlowStatsQueryResult is the result for a FlowStatQuery
type FlowStatsQueryResult struct {
	Node  UID
	Stats []FlowStats
}

// FlowStats is the statistics of flow
type FlowStats struct {
	Match    Match
	Duration time.Duration
	Packets  uint64
	Bytes    uint64
}

func (stats FlowStats) BW() Bandwidth {
	if stats.Duration == 0 {
		return 0
	}
	return Bandwidth(stats.Bytes / uint64(stats.Duration))
}

func init() {
	gob.Register(FlowStatsQuery{})
	gob.Register(FlowStatsQueryResult{})
	gob.Register(NodeQuery{})
	gob.Register(NodeQueryResult{})
	gob.Register(PortQuery{})
	gob.Register(PortQueryResult{})
}
