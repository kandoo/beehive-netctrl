package nom

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
)

// PortUpdated is a high-level event emitted when a port is added, removed, or
// its state/configuration is changed.
type PortUpdated Port

// PortStatusChanged is emitted when a driver receives a port status
type PortStatusChanged struct {
	Port   Port
	Driver Driver
}

// Port is either a physical or a virtual port of a node.
type Port struct {
	ID      PortID      // ID is unique among the ports of this node.
	Name    string      // Human-readable name of the port.
	MACAddr MACAddr     // Hardware address of the port.
	Node    UID         // The node.
	Link    UID         // The outgoing link.
	State   PortState   // Is the state of the port.
	Config  PortConfig  // Is the configuration of the port.
	Feature PortFeature // Features of this port.
}

func (p Port) String() string {
	return fmt.Sprintf("Port (node=%v, id=%v, addr=%v)", p.Node, p.ID, p.MACAddr)
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

// NodeFromPortUID returns the node UID from the port's UID.
func NodeFromPortUID(port UID) (node UID) {
	n, _ := ParsePortUID(port)
	return n.UID()
}

// JSONDecode decodes the port from a byte array using JSON.
func (p *Port) JSONDecode(b []byte) error {
	return json.Unmarshal(b, p)
}

// JSONEncode encodes the port into a byte array using JSON.
func (p *Port) JSONEncode() ([]byte, error) {
	return json.Marshal(p)
}

// Ports is a slice of ports with useful auxilaries.
type Ports []Port

// GetPort retrieves a port by its ID, and returns false if no port is found.
func (ports Ports) GetPort(id UID) (Port, bool) {
	for _, p := range ports {
		if p.UID() == id {
			return p, true
		}
	}
	return Port{}, false
}

// HasPort returns whether port is in ports.
func (ports Ports) HasPort(port Port) bool {
	_, ok := ports.GetPort(port.UID())
	return ok
}

// AddPort adds p to ports.
func (ports *Ports) AddPort(p Port) {
	*ports = append(*ports, p)
}

// DelPort deletes port from ports. If there is no such port, it returns false.
func (ports *Ports) DelPort(port Port) bool {
	id := port.UID()
	for i, p := range *ports {
		if p.UID() == id {
			*ports = append((*ports)[:i], (*ports)[i+1:]...)
			return true
		}
	}
	return false
}

// PortState is the current state of a port.
type PortState uint8

// Valid values for PortState.
const (
	PortStateUnknown PortState = iota // Port's state is unknown.
	PortStateDown              = iota // Port is not connected to any link.
	PortStateUp                = iota // Port is up and forwarding packets.
	PortStateBlocked           = iota // Port is blocked.
)

// PortConfig is the NOM specific configuration of the port.
type PortConfig uint8

// Valid values for PortConfig.
const (
	PortConfigDown        PortConfig = 1 << iota // Down.
	PortConfigDropPackets            = 1 << iota // Drop incoming packets.
	PortConfigNoForward              = 1 << iota // Do not forward packets.
	PortConfigNoFlood                = 1 << iota // Do not include in flood.
	PortConfigNoPacketIn             = 1 << iota // Do not send packet ins.
	PortConfigDisableStp             = 1 << iota // Disable STP.
	PortConfigDropStp                = 1 << iota // Drop STP packets.
)

// PortFeature represents port features.
type PortFeature uint16

// Valid values for PortFeature
const (
	PortFeature10MBHD  PortFeature = 1 << iota // 10MB half-duplex.
	PortFeature10MBFD              = 1 << iota // 10MB full-duplex.
	PortFeature100MBHD             = 1 << iota // 100MB half-duplex.
	PortFeature100MBFD             = 1 << iota // 100MB half-duplex.
	PortFeature1GBHD               = 1 << iota // 1GB half-duplex.
	PortFeature1GBFD               = 1 << iota // 1GB half-duplex.
	PortFeature10GBHD              = 1 << iota // 10GB  half-duplex.
	PortFeature10GBFD              = 1 << iota // 10GB half-duplex.
	PortFeatureCopper              = 1 << iota // Copper.
	PortFeatureFiber               = 1 << iota // Fiber.
	PortFeatureAutoneg             = 1 << iota // Auto negotiation.
	PortPause                      = 1 << iota // Pause.
	PortPauseAsym                  = 1 << iota // Asymmetric pause.
)

func init() {
	gob.Register(Port{})
	gob.Register(PortID(""))
	gob.Register(PortStatusChanged{})
	gob.Register(PortUpdated{})
}
