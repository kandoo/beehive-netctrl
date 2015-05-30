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

// JSONDecode decodes the network from a byte slice using JSON.
func (n *Network) JSONDecode(b []byte) error {
	return json.Unmarshal(b, n)
}

// JSONEncode encodes the network into a byte slice using JSON.
func (n *Network) JSONEncode() ([]byte, error) {
	return json.Marshal(n)
}
