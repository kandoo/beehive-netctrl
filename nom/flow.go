package nom

import (
	"encoding/gob"
	"time"

	bh "github.com/kandoo/beehive"
)

// InPort is the input port field.
type InPort UID

func (in InPort) HasSameType(f Field) bool {
	switch f.(type) {
	case InPort:
		return true
	}
	return false
}

func (in InPort) Equals(f Field) bool {
	switch field := f.(type) {
	case InPort:
		return in == field
	}
	return false
}

func (in InPort) Subsumes(f Field) bool {
	return in.Equals(f)
}

// EthAddrField is a common type for EthDst and EthSrc.
type EthAddrField struct {
	Addr MACAddr
	Mask MACAddr
}

// EthDst is the field for Ethernet destination address.
type EthDst MaskedMACAddr

func (e EthDst) HasSameType(f Field) bool {
	switch f.(type) {
	case EthDst:
		return true
	}
	return false
}

func (e EthDst) Equals(f Field) bool {
	switch field := f.(type) {
	case EthDst:
		return e == field
	}
	return false
}

func (e EthDst) Subsumes(f Field) bool {
	switch field := f.(type) {
	case EthDst:
		return MaskedMACAddr(e).Subsumes(MaskedMACAddr(field))
	}
	return false
}

// EthSrc is the field for Ethernet source address.
type EthSrc MaskedMACAddr

func (e EthSrc) HasSameType(f Field) bool {
	switch f.(type) {
	case EthSrc:
		return true
	}
	return false
}

func (e EthSrc) Equals(f Field) bool {
	switch field := f.(type) {
	case EthSrc:
		return e == field
	}
	return false
}

func (e EthSrc) Subsumes(f Field) bool {
	switch field := f.(type) {
	case EthSrc:
		return MaskedMACAddr(e).Subsumes(MaskedMACAddr(field))
	}
	return false
}

// EthType represents the field for an Ethernet type.
type EthType uint16

func (e EthType) HasSameType(f Field) bool {
	switch f.(type) {
	case EthType:
		return true
	}
	return false
}

func (e EthType) Equals(f Field) bool {
	switch field := f.(type) {
	case EthType:
		return e == field
	}
	return false
}

func (e EthType) Subsumes(f Field) bool {
	return e.Equals(f)
}

// VLANID represents the field for the VLAN ID.
type VLANID uint16

func (e VLANID) HasSameType(f Field) bool {
	switch f.(type) {
	case VLANID:
		return true
	}
	return false
}

func (e VLANID) Equals(f Field) bool {
	switch field := f.(type) {
	case VLANID:
		return e == field
	}
	return false
}

func (e VLANID) Subsumes(f Field) bool {
	return e.Equals(f)
}

// VLANPCP represents the field for the VLAN PCP.
type VLANPCP uint8

func (e VLANPCP) HasSameType(f Field) bool {
	switch f.(type) {
	case VLANPCP:
		return true
	}
	return false
}

func (e VLANPCP) Equals(f Field) bool {
	switch field := f.(type) {
	case VLANPCP:
		return e == field
	}
	return false
}

func (e VLANPCP) Subsumes(f Field) bool {
	return e.Equals(f)
}

type IPv4Src MaskedIPv4Addr

func (ip IPv4Src) HasSameType(f Field) bool {
	switch f.(type) {
	case IPv4Src:
		return true
	}
	return false
}

func (ip IPv4Src) Equals(f Field) bool {
	switch field := f.(type) {
	case IPv4Src:
		return ip == field
	}
	return false
}

func (ip IPv4Src) Subsumes(f Field) bool {
	switch field := f.(type) {
	case IPv4Src:
		return MaskedIPv4Addr(ip).Subsumes(MaskedIPv4Addr(field))
	}
	return false
}

type IPv4Dst MaskedIPv4Addr

func (ip IPv4Dst) HasSameType(f Field) bool {
	switch f.(type) {
	case IPv4Dst:
		return true
	}
	return false
}

func (ip IPv4Dst) Equals(f Field) bool {
	switch field := f.(type) {
	case IPv4Dst:
		return ip == field
	}
	return false
}

func (ip IPv4Dst) Subsumes(f Field) bool {
	switch field := f.(type) {
	case IPv4Dst:
		return MaskedIPv4Addr(ip).Subsumes(MaskedIPv4Addr(field))
	}
	return false
}

type IPv6Src MaskedIPv6Addr

func (ip IPv6Src) HasSameType(f Field) bool {
	switch f.(type) {
	case IPv6Src:
		return true
	}
	return false
}

func (ip IPv6Src) Equals(f Field) bool {
	switch field := f.(type) {
	case IPv6Src:
		return ip == field
	}
	return false
}

func (ip IPv6Src) Subsumes(f Field) bool {
	switch field := f.(type) {
	case IPv6Src:
		return MaskedIPv6Addr(ip).Subsumes(MaskedIPv6Addr(field))
	}
	return false
}

type IPv6Dst MaskedIPv6Addr

func (ip IPv6Dst) HasSameType(f Field) bool {
	switch f.(type) {
	case IPv6Dst:
		return true
	}
	return false
}

func (ip IPv6Dst) Equals(f Field) bool {
	switch field := f.(type) {
	case IPv6Dst:
		return ip == field
	}
	return false
}

func (ip IPv6Dst) Subsumes(f Field) bool {
	switch field := f.(type) {
	case IPv6Dst:
		return MaskedIPv6Addr(ip).Subsumes(MaskedIPv6Addr(field))
	}
	return false
}

type TransportPortSrc uint16

func (p TransportPortSrc) HasSameType(f Field) bool {
	switch f.(type) {
	case TransportPortSrc:
		return true
	}
	return false
}

func (p TransportPortSrc) Equals(f Field) bool {
	switch field := f.(type) {
	case TransportPortSrc:
		return p == field
	}
	return false
}

func (p TransportPortSrc) Subsumes(f Field) bool {
	return p.Equals(f)
}

type TransportPortDst uint16

func (p TransportPortDst) HasSameType(f Field) bool {
	switch f.(type) {
	case TransportPortDst:
		return true
	}
	return false
}

func (p TransportPortDst) Equals(f Field) bool {
	switch field := f.(type) {
	case TransportPortDst:
		return p == field
	}
	return false
}

func (p TransportPortDst) Subsumes(f Field) bool {
	return p.Equals(f)
}

// Valid values for EthType.
const (
	EthTypeIPv4 EthType = 0x0800
	EthTypeIPv6         = 0x86DD
	EthTypeARP          = 0x0806
)

type Field interface {
	HasSameType(f Field) bool
	Equals(f Field) bool
	Subsumes(f Field) bool
}

// Match is a collection of fields that will match packets.
type Match struct {
	Fields []Field
}

// Clone creates a copy of the match.
func (m Match) Clone() Match {
	clone := Match{
		Fields: make([]Field, len(m.Fields)),
	}
	copy(clone.Fields, m.Fields)
	return clone
}

func (m Match) InPort() (InPort, bool) {
	for _, f := range m.Fields {
		switch field := f.(type) {
		case InPort:
			return field, true
		}
	}
	return InPort(0), false
}

func (m Match) EthType() (EthType, bool) {
	for _, f := range m.Fields {
		switch field := f.(type) {
		case EthType:
			return field, true
		}
	}
	return EthType(0), false
}

func (m Match) IPv4Src() (IPv4Src, bool) {
	for _, f := range m.Fields {
		switch field := f.(type) {
		case IPv4Src:
			return field, true
		}
	}
	return IPv4Src{}, false
}

func (m Match) IPv4Dst() (IPv4Dst, bool) {
	for _, f := range m.Fields {
		switch field := f.(type) {
		case IPv4Dst:
			return field, true
		}
	}
	return IPv4Dst{}, false
}

func (m Match) IPv6Src() (IPv6Src, bool) {
	for _, f := range m.Fields {
		switch field := f.(type) {
		case IPv6Src:
			return field, true
		}
	}
	return IPv6Src{}, false
}

func (m Match) IPv6Dst() (IPv6Dst, bool) {
	for _, f := range m.Fields {
		switch field := f.(type) {
		case IPv6Dst:
			return field, true
		}
	}
	return IPv6Dst{}, false
}

func (m Match) VLANID() (VLANID, bool) {
	for _, f := range m.Fields {
		switch field := f.(type) {
		case VLANID:
			return field, true
		}
	}
	return VLANID(0), false
}

func (m Match) VLANPCP() (VLANPCP, bool) {
	for _, f := range m.Fields {
		switch field := f.(type) {
		case VLANPCP:
			return field, true
		}
	}
	return VLANPCP(0), false
}

func (m Match) TransportPortSrc() (TransportPortSrc, bool) {
	for _, f := range m.Fields {
		switch field := f.(type) {
		case TransportPortSrc:
			return field, true
		}
	}
	return TransportPortSrc(0), false
}

func (m Match) TransportPortDst() (TransportPortDst, bool) {
	for _, f := range m.Fields {
		switch field := f.(type) {
		case TransportPortDst:
			return field, true
		}
	}
	return TransportPortDst(0), false
}

func (m Match) Subsumes(thatm Match) bool {
	for _, thisf := range m.Fields {
		if thatm.countFields(thisf.HasSameType) != 1 {
			return false
		}
		if thatm.countFields(thisf.Subsumes) != 1 {
			return false
		}
	}
	return true
}

func (m Match) Equals(thatm Match) bool {
	for _, thisf := range m.Fields {
		if thatm.countFields(thisf.HasSameType) != 1 {
			return false
		}
		if thatm.countFields(thisf.Equals) != 1 {
			return false
		}
	}
	return true
}

func (m Match) assertFields(checker func(f Field) bool) bool {
	for _, f := range m.Fields {
		if !checker(f) {
			return false
		}
	}
	return true
}

func (m Match) countFields(checker func(f Field) bool) int {
	count := 0
	for _, f := range m.Fields {
		if checker(f) {
			count++
		}
	}
	return count
}

type Action interface {
	Equals(a Action) bool
}

type ActionForward struct {
	Ports []UID
}

func (a ActionForward) Equals(thata Action) bool {
	thataf, ok := thata.(ActionForward)
	if !ok {
		return false
	}
	ports := make(map[UID]struct{})
	for _, p := range a.Ports {
		ports[p] = struct{}{}
	}
	for _, p := range thataf.Ports {
		if _, ok := ports[p]; !ok {
			return false
		}
	}
	return true
}

type ActionDrop struct {
}

func (a ActionDrop) Equals(thata Action) bool {
	_, ok := thata.(ActionDrop)
	return ok
}

type ActionFlood struct {
	InPort UID
}

func (a ActionFlood) Equals(thata Action) bool {
	thataf, ok := thata.(ActionFlood)
	if !ok {
		return false
	}
	return a.InPort == thataf.InPort
}

type ActionSendToController struct{}

func (a ActionSendToController) Equals(thata Action) bool {
	_, ok := thata.(ActionSendToController)
	return ok
}

type ActionPushVLAN struct {
	ID VLANID
}

func (a ActionPushVLAN) Equals(thata Action) bool {
	thatap, ok := thata.(ActionPushVLAN)
	if !ok {
		return false
	}
	return a.ID == thatap.ID
}

type ActionPopVLAN struct{}

func (a ActionPopVLAN) Equals(thata Action) bool {
	_, ok := thata.(ActionPopVLAN)
	return ok
}

type ActionWriteFields struct {
	Fields []Field
}

func (a ActionWriteFields) Equals(thata Action) bool {
	thataw, ok := thata.(ActionWriteFields)
	if !ok {
		return false
	}
	for i := range a.Fields {
		if !a.Fields[i].Equals(thataw.Fields[i]) {
			return false
		}
	}
	return true
}

// FlowEntry represents a match-action rule for a specific node.
type FlowEntry struct {
	ID          string // ID is defined by the subscriber, not globally unique.
	Node        UID
	Match       Match
	Actions     []Action
	Priority    uint16
	IdleTimeout time.Duration
	HardTimeout time.Duration
}

func (f FlowEntry) Equals(thatf FlowEntry) bool {
	for i := range f.Actions {
		if f.Actions[i].Equals(thatf.Actions[i]) {
			return false
		}
	}
	return f.Node == thatf.Node && f.Match.Equals(thatf.Match) &&
		f.Priority == thatf.Priority
}

// AddFlowEntry is a message emitted to install a flow entry on a node.
type AddFlowEntry struct {
	Subscriber bh.AppCellKey
	Flow       FlowEntry
}

// DelFlowEntry is emitted to remove the flow entries with the given match.
// If Exact is false, it removes all flow entries that are subsumed by the
// given match.
type DelFlowEntry struct {
	Match Match
	Exact bool
}

// FlowEntryDeleted is emitted (broadcasted and also sent to the subscriber of
// the flow) when a flow is deleted.
type FlowEntryDeleted struct {
	Flow FlowEntry
}

// FlowEntryAdded is emitted (broadcasted and also sent to the subscriber of the
// flow) when a flow is added. If the flow already existed, the message is
// emitted to the subscriber.
type FlowEntryAdded struct {
	Flow FlowEntry
}

func init() {
	gob.Register(ActionDrop{})
	gob.Register(ActionFlood{})
	gob.Register(ActionForward{})
	gob.Register(ActionPopVLAN{})
	gob.Register(ActionPushVLAN{})
	gob.Register(ActionSendToController{})
	gob.Register(ActionWriteFields{})
	gob.Register(AddFlowEntry{})
	gob.Register(DelFlowEntry{})
	gob.Register(EthAddrField{})
	gob.Register(EthDst{})
	gob.Register(EthSrc{})
	gob.Register(EthType(0))
	gob.Register(FlowEntryAdded{})
	gob.Register(FlowEntryDeleted{})
	gob.Register(FlowEntry{})
	gob.Register(IPv4Dst{})
	gob.Register(IPv4Src{})
	gob.Register(IPv6Dst{})
	gob.Register(IPv6Src{})
	gob.Register(InPort(0))
	gob.Register(Match{})
	gob.Register(TransportPortDst(0))
	gob.Register(TransportPortSrc(0))
	gob.Register(VLANID(0))
	gob.Register(VLANPCP(0))
}
