package nom

import (
	"encoding/gob"
	"fmt"
	"time"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive/strings"
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

func (in InPort) String() string {
	return fmt.Sprintf("in_port=%v", UID(in))
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

func (e EthDst) String() string {
	return fmt.Sprintf("eth_dst=%v", MaskedMACAddr(e))
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

func (e EthSrc) String() string {
	return fmt.Sprintf("eth_src=%v", MaskedMACAddr(e))
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

func (e EthType) String() string {
	return fmt.Sprintf("eth_type=%v", uint16(e))
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

func (e VLANID) String() string {
	return fmt.Sprintf("vlan=%v", uint16(e))
}

// VLANPCP represents the field for the VLAN PCP.
type VLANPCP uint8

func (p VLANPCP) HasSameType(f Field) bool {
	_, ok := f.(VLANPCP)
	return ok
}

func (p VLANPCP) Equals(f Field) bool {
	if fp, ok := f.(VLANPCP); ok {
		return p == fp
	}
	return false
}

func (p VLANPCP) Subsumes(f Field) bool {
	return p.Equals(f)
}

func (p VLANPCP) String() string {
	return fmt.Sprintf("vlan_pcp=%v", uint8(p))
}

type IPProto uint8

func (p IPProto) HasSameType(f Field) bool {
	_, ok := f.(IPProto)
	return ok
}

func (p IPProto) Equals(f Field) bool {
	if fp, ok := f.(IPProto); ok {
		return p == fp
	}
	return false
}

func (p IPProto) Subsumes(f Field) bool {
	return p.Equals(f)
}

func (p IPProto) String() string {
	return fmt.Sprint("ip_proto=%v", uint8(p))
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

func (ip IPv4Src) String() string {
	return MaskedIPv4Addr(ip).String()
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

func (ip IPv4Dst) String() string {
	return fmt.Sprintf("ipv4_dst=%v", MaskedIPv4Addr(ip))
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

func (ip IPv6Src) String() string {
	return fmt.Sprintf("ipv6_src=%v", MaskedIPv6Addr(ip))
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

func (ip IPv6Dst) String() string {
	return fmt.Sprintf("ipv6_dst=%v", MaskedIPv6Addr(ip))
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

func (p TransportPortSrc) String() string {
	return fmt.Sprintf("tp_port_src=%v", uint16(p))
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

func (p TransportPortDst) String() string {
	return fmt.Sprintf("tp_port_dst=%v", uint16(p))
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

func (m Match) String() string {
	l := len(m.Fields)
	if l == 0 {
		return "match(*)"
	}
	a := make([]interface{}, len(m.Fields))
	for i := range m.Fields {
		a[i] = m.Fields[i]
	}
	return fmt.Sprintf("match(%v)", strings.Join(a, ","))
}

// AddField adds the field to the list of fields in the match.
func (m *Match) AddField(f Field) {
	m.Fields = append(m.Fields, f)
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

func (m Match) EthSrc() (EthSrc, bool) {
	for _, f := range m.Fields {
		switch field := f.(type) {
		case EthSrc:
			return field, true
		}
	}
	return EthSrc{}, false
}

func (m Match) EthDst() (EthDst, bool) {
	for _, f := range m.Fields {
		switch field := f.(type) {
		case EthDst:
			return field, true
		}
	}
	return EthDst{}, false
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
	if len(m.Fields) != len(thatm.Fields) {
		return false
	}
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

func (a ActionForward) String() string {
	return fmt.Sprintf("forward(to=%v)", a.Ports)
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

func (a ActionDrop) String() string {
	return fmt.Sprintf("drop")
}

type ActionFlood struct {
	InPort UID
}

func (a ActionFlood) String() string {
	return fmt.Sprintf("flood(except=%v)", a.InPort)
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

func (f FlowEntry) String() string {
	a := make([]interface{}, len(f.Actions))
	for i := range f.Actions {
		a[i] = f.Actions[i]
	}
	astr := strings.Join(a, ",")
	return fmt.Sprintf("flow(%v=>%v,node=%v,priority=%v,idleto=%v,hardto=%v)",
		f.Match, astr, f.Node, f.Priority, f.IdleTimeout, f.HardTimeout)
}

func (f FlowEntry) Equals(thatf FlowEntry) bool {
	if f.Node != thatf.Node || f.Priority != thatf.Priority ||
		len(f.Actions) != len(thatf.Actions) {

		return false
	}
	for i := range f.Actions {
		if !f.Actions[i].Equals(thatf.Actions[i]) {
			return false
		}
	}
	return f.Match.Equals(thatf.Match)
}

// Subsumes returns whether everything in f is equal to thatf except that f's
// match subsumes thatf's match.
func (f FlowEntry) Subsumes(thatf FlowEntry) bool {
	if f.Node != thatf.Node || f.Priority != thatf.Priority ||
		len(f.Actions) != len(thatf.Actions) {

		return false
	}
	for i := range f.Actions {
		if !f.Actions[i].Equals(thatf.Actions[i]) {
			return false
		}
	}
	return f.Match.Subsumes(thatf.Match)
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
