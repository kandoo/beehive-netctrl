package nom

import (
	"encoding/gob"
	"encoding/json"
)

// LinkAdded is emitted when a new link is discovered.
type LinkAdded Link

// LinkRemoved is emitted when a new link is removed.
type LinkRemoved Link

// Link represents an outgoing link from a port.
type Link struct {
	ID    LinkID    // Link's ID.
	From  UID       // From is the link's port.
	To    []UID     // To stores the port(s) connected to From using this link.
	State LinkState // The link's state.
}

// LinkID is a link's ID which is unique among the outgoing links of a port.
type LinkID string

// UID returns the UID of the link in the form of
// net_id$$node_id$$port_id$$link_id.
func (l Link) UID() UID {
	return UIDJoin(string(l.From), string(l.ID))
}

// ParseLinkUID parses a link UID into the respetive node, port, and link ids.
func ParseLinkUID(id UID) (NodeID, PortID, LinkID) {
	s := UIDSplit(id)
	return NodeID(s[0]), PortID(s[1]), LinkID(s[2])
}

// GoDecode decodes the link from b using Gob.
func (l *Link) GoDecode(b []byte) error {
	return ObjGoDecode(l, b)
}

// GoEncode encodes the node into a byte array using Gob.
func (l *Link) GoEncode() ([]byte, error) {
	return ObjGoEncode(l)
}

// JSONDecode decodes the node from a byte array using JSON.
func (l *Link) JSONDecode(b []byte) error {
	return json.Unmarshal(b, l)
}

// JSONEncode encodes the node into a byte array using JSON.
func (l *Link) JSONEncode() ([]byte, error) {
	return json.Marshal(l)
}

// LinkState represents the status of a link.
type LinkState uint8

// Valid values for LinkState.
const (
	LinkStateUnknown LinkState = iota
	LinkStateUp                = iota
	LinkStateDown              = iota
)

// Bandwidth represents bandwidth in Bps.
type Bandwidth uint64

// Bandwidth units.
const (
	KBps Bandwidth = 1000
	MBps Bandwidth = 1000000
	GBps Bandwidth = 1000000000
)

func init() {
	gob.Register(LinkAdded{})
	gob.Register(LinkRemoved{})
	gob.Register(Link{})
	gob.Register(LinkID(""))
	gob.Register(LinkState(0))
	gob.Register(Bandwidth(0))
}
