package nom

// Path represents a directed, logical multi-path among points in the network.
// Each path start from one or more points in the network and is then connected
// to other paths.
type Path struct {
	ID UID // A universally unique ID of this path.

	Match Match          // Match matches the flows that go into this path.
	Src   PointGroup     // Src is the points that start this path.
	Dst   []WeightedPath // Dst is the destination(s) of this path.

	Redundancy int       // Number of redundant links to dst.
	MinBW      Bandwidth // Minimum required BW of this path.
	MaxBW      Bandwidth // Maximum required BW of this path.
}

// WeightedPath represents a path along with its weight. This is only used as a
// destination for other paths.
type WeightedPath struct {
	Weight int
	Path   *Path
}

// Point represents a connecting point in a path, denoted by an input port and
// an output port. Input and output ports must be of the same node.
//
// The starting points of a path have no input port, and the terminal points
// have no output port.
//
// For intermediary points, if the output port is "" meaning that the output
// port should be automatically selected based on the next hub. If the input
// port is "", it means that flows can be received from all ports of that node.
type Point struct {
	Node UID // The node.
	In   UID // Input port.
	Out  UID // Output port.
}

// PointGroup represents a collection of points in a path.
type PointGroup struct {
	Points []Point
}

// NewPortGroup creates a group of path points.
func NewPointGroup(p ...Point) *PointGroup {
	return &PointGroup{
		Points: p,
	}
}

// CreatePath is a message emitted to create a path in the network.
type CreatePath Path

// DeletePath is a message emitted to delete a path.
type DeletePath struct {
	Path UID
}

// ReplacePath is a message emitted to replace an old path with a new one.
type ReplacePath struct {
	OldPath Path
	NewPath Path
}

// QueryPath is emitted to query the stats of a path.
type QueryPath struct {
	Path UID
}

// PathStatistics is emitted as a reply to a QueryPath.
type QueryPathResult struct {
	Path  UID
	Stats []PointStats
}

// PointStats is the statistics of a specific point.
type PointStats struct {
	Point   Point
	Bytes   uint64
	Packets uint64
}
