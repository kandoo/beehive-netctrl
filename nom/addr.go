package nom

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// MACAddr represents a MAC address.
type MACAddr [6]byte

var (
	MaskNoneMAC            MACAddr   = [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	BroadcastMAC           MACAddr   = [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	CDPMulticastMAC        MACAddr   = [6]byte{0x01, 0x00, 0x0C, 0xCC, 0xCC, 0xCC}
	CiscoSTPMulticastMAC   MACAddr   = [6]byte{0x01, 0x00, 0x0C, 0xCC, 0xCC, 0xCD}
	IEEE802MulticastPrefix MACAddr   = [6]byte{0x01, 0x80, 0xC2, 0x00, 0x00, 0x00}
	IPv4MulticastPrefix    MACAddr   = [6]byte{0x01, 0x00, 0x5E, 0x00, 0x00, 0x00}
	IPv6MulticastPrefix    MACAddr   = [6]byte{0x33, 0x33, 0x00, 0x00, 0x00, 0x00}
	LLDPMulticastMACs      []MACAddr = []MACAddr{
		[6]byte{0x01, 0x80, 0xC2, 0x00, 0x00, 0x0E},
		[6]byte{0x01, 0x80, 0xC2, 0x00, 0x00, 0x0C},
		[6]byte{0x01, 0x80, 0xC2, 0x00, 0x00, 0x00},
	}
	MaskNoneIPV4 IPv4Addr = [4]byte{0xFF, 0xFF, 0xFF, 0xFF}
	MaskNoneIPV6 IPv6Addr = [16]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
)

func (m MACAddr) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		m[0], m[1], m[2], m[3], m[4], m[5])
}

// Key returns an string represtation of the MAC address suitable to store in
// dictionaries. It is more efficient compared to MACAddr.String().
func (m MACAddr) Key() string {
	return string(m[:])
}

// IsBroadcast returns whether the MAC address is a broadcast address.
func (m MACAddr) IsBroadcast() bool {
	return m == BroadcastMAC
}

// IsMulticast returns whether the MAC address is a multicast address.
func (m MACAddr) IsMulticast() bool {
	return m == CDPMulticastMAC || m == CiscoSTPMulticastMAC ||
		m.hasPrefix(IEEE802MulticastPrefix, 3) ||
		m.hasPrefix(IPv4MulticastPrefix, 3) ||
		m.hasPrefix(IPv6MulticastPrefix, 2)
}

// IsLLDP returns whether the mac address is a multicast address used for LLDP.
func (m MACAddr) IsLLDP() bool {
	for _, lm := range LLDPMulticastMACs {
		if m == lm {
			return true
		}
	}
	return false
}

func (m MACAddr) hasPrefix(p MACAddr, l int) bool {
	for i := 0; i < l; i++ {
		if p[i] != m[i] {
			return false
		}
	}
	return true
}

func (m MACAddr) Mask(mask MACAddr) MACAddr {
	masked := m
	for i := range mask {
		masked[i] &= mask[i]
	}
	return masked
}

func (m MACAddr) Less(thatm MACAddr) bool {
	for i := range thatm {
		switch {
		case m[i] < thatm[i]:
			return true
		case m[i] > thatm[i]:
			return false
		}
	}
	return false
}

// MaskedMACAddr is a MAC address that is wildcarded with a mask.
type MaskedMACAddr struct {
	Addr MACAddr // The MAC address.
	Mask MACAddr // The mask of the MAC address.
}

// Match returns whether the masked mac address matches mac.
func (mm MaskedMACAddr) Match(mac MACAddr) bool {
	return mm.Mask.Mask(mm.Addr) == mm.Mask.Mask(mac)
}

// Subsumes returns whether this mask address includes all the addresses matched
// by thatmm.
func (mm MaskedMACAddr) Subsumes(thatmm MaskedMACAddr) bool {
	if thatmm.Mask.Less(mm.Mask) {
		return false
	}
	return mm.Match(thatmm.Addr.Mask(thatmm.Mask))
}

// IPv4Addr represents an IP version 4 address in big endian byte order.
// For example, 127.0.0.1 is represented as IPv4Addr{127, 0, 0, 1}.
type IPv4Addr [4]byte

// Uint converts the IP version 4 address into a 32-bit integer in little
// endian byte order.
func (ip IPv4Addr) Uint() uint32 {
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 |
		uint32(ip[3])
}

// FromUint loads the ip address from addr.
func (ip *IPv4Addr) FromUint(addr uint32) {
	for i := 0; i < 4; i++ {
		ip[i] = byte((addr >> uint(8*(3-i))) & 0xFF)
	}
}

// PopCount returns the number of ones in the IP address. For example, it
// returns 16 for IPv4Addr{255, 255, 0, 0}.
func (ip IPv4Addr) PopCount() uint32 {
	v := ip.Uint()
	v -= (v >> 1) & 0x55555555
	v = ((v >> 2) & 0x33333333) + v&0x33333333
	v = ((v >> 4) + v) & 0x0F0F0F0F
	v = ((v >> 8) + v) & 0x00FF00FF
	return ((v >> 16) + v) & 0x0000FFFF
}

// Mask masked the IP address with mask.
func (ip IPv4Addr) Mask(mask IPv4Addr) IPv4Addr {
	masked := ip
	for i := range masked {
		masked[i] &= mask[i]
	}
	return masked
}

// Less returns whether ip is less than thatip.
func (ip IPv4Addr) Less(thatip IPv4Addr) bool {
	for i := range ip {
		switch {
		case ip[i] < thatip[i]:
			return true
		case ip[i] > thatip[i]:
			return false
		}
	}
	return false
}

// AsCIDRMask returns the CIDR prefix number based on this address.
// For example, it returns 24 for 255.255.255.0.
func (ip IPv4Addr) AsCIDRMask() int {
	m := 0
	for i := len(ip) - 1; i >= 0; i-- {
		for j := uint(0); j < 8; j++ {
			if (ip[i]>>j)&0x1 != 0 {
				return 32 - m
			}
			m++
		}
	}
	return 0
}

func (ip IPv4Addr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

// CIDRToMaskedIPv4 converts a CIDR-style IP address into a NOM masked IP
// address. For example, if addr is 0x7F000001 and mask is 24, this function
// returns {IPv4Addr{127, 0, 0, 1}, IPv4Addr{255, 255, 255, 0}}.
func CIDRToMaskedIPv4(addr uint32, mask uint) MaskedIPv4Addr {
	var maskedip MaskedIPv4Addr
	maskedip.Addr.FromUint(addr)
	maskedip.Mask.FromUint(uint32(((1 << mask) - 1) << (32 - mask)))
	return maskedip
}

// MaskedIPv4Addr represents a masked IP address (ie, an IPv4 prefix)
type MaskedIPv4Addr struct {
	Addr IPv4Addr
	Mask IPv4Addr
}

// Match returns whether the masked IP address matches ip.
func (mi MaskedIPv4Addr) Match(ip IPv4Addr) bool {
	return mi.Addr.Mask(mi.Mask) == ip.Mask(mi.Mask)
}

func (mi MaskedIPv4Addr) Subsumes(thatmi MaskedIPv4Addr) bool {
	if thatmi.Mask.Less(mi.Mask) {
		return false
	}
	return mi.Addr.Mask(mi.Mask) == thatmi.Addr.Mask(mi.Mask)
}

func (mi MaskedIPv4Addr) String() string {
	return fmt.Sprintf("%v/%d", mi.Addr, mi.Mask.AsCIDRMask())
}

// IPv6Addr represents an IP version 6 address in big-endian byte order.
type IPv6Addr [16]byte

// Mask masked the IP address with mask.
func (ip IPv6Addr) Mask(mask IPv6Addr) IPv6Addr {
	masked := ip
	for i := range masked {
		masked[i] &= mask[i]
	}
	return masked
}

// Less returns whether ip is less than thatip.
func (ip IPv6Addr) Less(thatip IPv6Addr) bool {
	for i := range ip {
		switch {
		case ip[i] < thatip[i]:
			return true
		case ip[i] > thatip[i]:
			return false
		}
	}
	return false
}

func (ip IPv6Addr) String() string {
	var buf bytes.Buffer
	zeros := 0
	for i := 0; i < 16; i += 2 {
		if ip[i] == 0 && ip[i+1] == 0 {
			zeros++
			if zeros == 2 {
				if zeros == i {
					buf.WriteString("::")
				} else {
					buf.WriteString(":")
				}
			}
			continue
		}
		if zeros == 1 {
			buf.WriteString("0:")
		}
		buf.WriteString(fmt.Sprintf("%x", int(ip[i])<<8|int(ip[i+1])))
		if i != 14 {
			buf.WriteString(":")
		}
		zeros = 0
	}
	return buf.String()
}

// AsCIDRMask returns the CIDR prefix number based on this address.
func (ip IPv6Addr) AsCIDRMask() int {
	m := 0
	for i := len(ip) - 1; i >= 0; i-- {
		for j := uint(0); j < 8; j++ {
			if (ip[i]>>j)&0x1 != 0 {
				return 128 - m
			}
			m++
		}
	}
	return 0
}

// MaskedIPv6Addr represents a masked IPv6 address.
type MaskedIPv6Addr struct {
	Addr IPv6Addr
	Mask IPv6Addr
}

// Match returns whether the masked IP address matches ip.
func (mi MaskedIPv6Addr) Match(ip IPv6Addr) bool {
	return mi.Addr.Mask(mi.Mask) == ip.Mask(mi.Mask)
}

func (mi MaskedIPv6Addr) Subsumes(thatmi MaskedIPv6Addr) bool {
	if thatmi.Mask.Less(mi.Mask) {
		return false
	}
	return mi.Addr.Mask(mi.Mask) == thatmi.Addr.Mask(mi.Mask)
}

func (mi MaskedIPv6Addr) String() string {
	return fmt.Sprintf("%v/%d", mi.Addr, mi.Mask.AsCIDRMask())
}

func init() {
	gob.Register(IPv4Addr{})
	gob.Register(IPv6Addr{})
	gob.Register(MACAddr{})
	gob.Register(MaskedIPv4Addr{})
	gob.Register(MaskedIPv6Addr{})
	gob.Register(MaskedMACAddr{})
}
