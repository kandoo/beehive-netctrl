package nom

import "fmt"

// MACAddr represents a MAC address.
type MACAddr [6]byte

func (m MACAddr) String() string {
	return fmt.Sprintf("%x:%x:%x:%x:%x:%x", m[0], m[1], m[2], m[3], m[4], m[5])
}

// IPV4Addr represents an IP version 4 address.
type IPV4Addr [4]byte

// IPV6Addr represents an IP version 6 address.
type IPV6Addr [16]byte
