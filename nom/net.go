package nom

import "encoding/json"

// Network represents a virtual or physical network in NOM.
type Network struct {
	ID   NetworkID // The id of the network.
	Desc string    // A human-readable description of the network.
}

// NetworkID is the ID of the network.
type NetworkID string

// UID returns the UID of the network.
func (n Network) UID() UID {
	return UID(n.ID)
}

// ParseNetworkUID parses a network UID into the network ID. Note that for
// network these IDs of the same.
func ParseNetworkUID(id UID) NetworkID {
	return NetworkID(id)
}

// GobDecode decodes the network from the byte slice using Gob.
func (n *Network) GobDecode(b []byte) error {
	return ObjGobDecode(n, b)
}

// GobEncode encodes the network into a byte slice using Gob.
func (n *Network) GobEncode() ([]byte, error) {
	return ObjGobEncode(n)
}

// JSONDecode decodes the network from a byte slice using JSON.
func (n *Network) JSONDecode(b []byte) error {
	return json.Unmarshal(b, n)
}

// JSONEncode encodes the network into a byte slice using JSON.
func (n *Network) JSONEncode() ([]byte, error) {
	return json.Marshal(n)
}
