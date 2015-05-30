package nom

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
)

// NodeConnected is a message emitted when a node connects to a driver.
type NodeConnected struct {
	Node   Node
	Driver Driver
}

// NodeDisconnected is a message emitted when a node disconnects from its
// driver.
type NodeDisconnected struct {
	Node   Node
	Driver Driver
}

// NodeJoined is a message emitted when a node joins the network through the
// controller. It is always emitted after processing NodeConnected in the
// controller.
type NodeJoined Node

// NodeLeft is a message emitted when a node disconnects from its driver. It is
// always emitted after processing NodeDisconnected in the controller.
type NodeLeft Node

// Node represents a forwarding element, such as switches and routers.
type Node struct {
	ID           NodeID
	Net          UID
	Capabilities []NodeCapability
	MACAddr      MACAddr
}

func (n Node) String() string {
	return fmt.Sprintf("Node %s (mac=%v)", string(n.ID), n.MACAddr)
}

// NodeID is the ID of a node. This must be unique among all nodes in the
// network.
type NodeID string

// UID converts id into a UID.
func (id NodeID) UID() UID {
	return UID(id)
}

// UID returns the node's unique ID. This id is in the form of net_id$$node_id.
func (n Node) UID() UID {
	return UID(string(n.ID))
}

// ParseNodeUID parses a UID of a node and returns the respective node IDs.
func ParseNodeUID(id UID) NodeID {
	s := UIDSplit(id)
	return NodeID(s[0])
}

// JSONDecode decodes the node from a byte array using JSON.
func (n *Node) JSONDecode(b []byte) error {
	return json.Unmarshal(b, n)
}

// JSONEncode encodes the node into a byte array using JSON.
func (n *Node) JSONEncode() ([]byte, error) {
	return json.Marshal(n)
}

func (n Node) HasCapability(c NodeCapability) bool {
	for _, nc := range n.Capabilities {
		if c == nc {
			return true
		}
	}

	return false
}

// NodeCapability is a capability of a NOM node.
type NodeCapability uint32

// Valid values for NodeCapability.
const (
	CapDriverRole NodeCapability = 1 << iota // Node can set the driver's role.
)

func init() {
	gob.Register(NodeConnected{})
	gob.Register(NodeDisconnected{})
	gob.Register(NodeLeft{})
	gob.Register(NodeJoined{})
}
