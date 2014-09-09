package nom

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
)

// Port is either a physical or a virtual port of a node.
type Port struct {
	ID    PortID
	Node  UID
	Links []UID
}

// PortID is the ID of a port and is unique among the ports of a node.
type PortID string

// UID returns the unique ID of the port in the form of
// net_id$$node_id$$port_id.
func (p Port) UID() UID {
	return UIDJoin(string(p.Node), string(p.ID))
}

// ParsePortUID parses a UID of a port and returns the respective node and port
// IDs.
func ParsePortUID(id UID) (NodeID, PortID) {
	s := UIDSplit(id)
	return NodeID(s[0]), PortID(s[1])
}

// GobDecode decodes the port from b using Gob.
func (p *Port) GobDecode(b []byte) error {
	return ObjGobDecode(p, b)
}

// GobEncode encodes the port into a byte array using Gob.
func (p *Port) GobEncode() ([]byte, error) {
	return ObjGobEncode(p)
}

// JSONDecode decodes the port from a byte array using JSON.
func (p *Port) JSONDecode(b []byte) error {
	return json.Unmarshal(b, p)
}

// JSONEncode encodes the port into a byte array using JSON.
func (p *Port) JSONEncode() ([]byte, error) {
	return json.Marshal(p)
}
