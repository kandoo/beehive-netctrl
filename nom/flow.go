package nom

import "time"

// InPort is the input port field.
type InPort UID

// EthDst is the field for Ethernet destination address.
type EthDst struct {
	Addr MACAddr
	Mask MACAddr
}

// EthSrc is the field for Ethernet source address.
type EthSrc struct {
	Addr MACAddr
	Mask MACAddr
}

// EthType represents the field for an Ethernet type.
type EthType uint16

// VLANID represents the field for the VLAN ID.
type VLANID uint16

// VLANPCP represents the field for the VLAN PCP.
type VLANPCP uint8

type IPV4Src struct {
	Addr IPV4Addr
	Mask IPV4Addr
}

type IPV4Dst struct {
	Addr IPV4Addr
	Mask IPV4Addr
}

type IPV6Src struct {
	Addr IPV6Addr
	Mask IPV6Addr
}

type IPV6Dst struct {
	Addr IPV6Addr
	Mask IPV6Addr
}

type TransportPortSrc uint16
type TransportPortDst uint16

// Valid values for EthType.
const (
	EthTypeIPV4 EthType = 0x0800
	EthTypeIPV6         = 0x86DD
	EthTypeARP          = 0x0806
)

type Field interface{}

type Match struct {
	Fields []Field
}

type Action interface{}

type ActionForward struct {
	Ports []UID
}

type ActionDrop struct {
}

type ActionFlood struct {
	InPort UID
}

type ActionSendToController struct {
}

type ActionPushVLAN struct {
	ID VLANID
}

type ActionPopVLAN struct {
}

type WriteFields struct {
	Fields []Field
}

// FlowPriority is the priority of a flow.
type FlowPriority uint16

// FlowEntry represents a match-action rule for a specific node.
type FlowEntry struct {
	Node        UID
	Match       Match
	Actions     []Action
	Priority    FlowPriority
	IdleTimeout time.Duration
	HardTimeout time.Duration
}

// AddFlowEntry is emitted to install a flow entry.
type AddFlowEntry struct {
	Flow FlowEntry
}

// DelFlowEntry is emitted to remove the flow entries with the given match.
// If Exact is false, it removes all flow entries that are subsumed by the
// given match.
type DelFlowEntry struct {
	Match Match
	Exact bool
}

// AddFlowEntryResult is emitted in response to a AddFlowEntry.
type AddFlowEntryResult struct {
	Err error
	Add AddFlowEntry
}

// DelFlowEntryResult is emitted in response to a DelFlowEntry.
type DeldFlowEntryResult struct {
	Err error
	Del DelFlowEntry
}
