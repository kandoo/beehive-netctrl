package nom

import "encoding/json"

// PortAdded is emitted when a port is added to a node (or when the node
// joins the network for the first time).
type PortAdded Port

// PortRemoved is emitted when a port is removed (or its node is disconnected
// from the controller).
type PortRemoved Port

// PortChanged is emitted when a port's state or configuration is changed.
type PortChanged Port

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
