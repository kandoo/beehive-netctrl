package nom

import "fmt"

// MACAddr represents a MAC address.
type MACAddr [6]byte

var (
	BroadcastMAC           MACAddr = [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	CDPMulticastMAC        MACAddr = [6]byte{0x01, 0x00, 0x0C, 0xCC, 0xCC, 0xCC}
	CiscoSTPMulticastMAC   MACAddr = [6]byte{0x01, 0x00, 0x0C, 0xCC, 0xCC, 0xCD}
	IEEE802MulticastPrefix MACAddr = [6]byte{0x01, 0x80, 0xC2, 0x00, 0x00, 0x00}
	IPV4MulticastPrefix    MACAddr = [6]byte{0x01, 0x00, 0x5E, 0x00, 0x00, 0x00}
	IPV6MulticastPrefix    MACAddr = [6]byte{0x33, 0x33, 0x00, 0x00, 0x00, 0x00}
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
		m.hasPrefix(IPV4MulticastPrefix, 3) ||
		m.hasPrefix(IPV6MulticastPrefix, 2)
}

func (m MACAddr) hasPrefix(p MACAddr, l int) bool {
	for i := 0; i < l; i++ {
		if p[i] != m[i] {
			return false
		}
	}
	return true
}

// IPV4Addr represents an IP version 4 address in little endian byte order.
// For example, 127.0.0.1 is represented as [4]byte{1, 0, 0, 127}.
type IPV4Addr [4]byte

func (ip IPV4Addr) Uint32() uint32 {
	return uint32(ip[0] | ip[1]<<8 | ip[2]<<16 | ip[3]<<24)
}

// IPV6Addr represents an IP version 6 address in little endian byte order.
type IPV6Addr [16]byte
