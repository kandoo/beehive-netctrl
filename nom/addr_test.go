package nom

import "testing"

func TestIPv4String(t *testing.T) {
	addrs := map[IPv4Addr]string{
		IPv4Addr{0, 0, 0, 0}:         "0.0.0.0",
		IPv4Addr{127, 0, 0, 1}:       "127.0.0.1",
		IPv4Addr{255, 255, 255, 255}: "255.255.255.255",
	}
	for a, s := range addrs {
		if a.String() != s {
			t.Errorf("invalid string for ipv4 address: actual=%v want=%v",
				a.String(), s)
		}
	}
}

func TestIPv4CIDR(t *testing.T) {
	addrs := map[IPv4Addr]int{
		IPv4Addr{0, 0, 0, 0}:         0,
		IPv4Addr{255, 128, 0, 0}:     9,
		IPv4Addr{255, 255, 0, 0}:     16,
		IPv4Addr{255, 255, 255, 255}: 32,
	}
	for a, m := range addrs {
		if a.AsCIDRMask() != m {
			t.Errorf("invalid CIDR for ipv4 mask: actual=%v want=%v",
				a.AsCIDRMask(), m)
		}
	}
}

func TestMaskedIPv4String(t *testing.T) {
	addrs := map[MaskedIPv4Addr]string{
		MaskedIPv4Addr{
			Addr: IPv4Addr{0, 0, 0, 0},
			Mask: IPv4Addr{},
		}: "0.0.0.0/0",
		MaskedIPv4Addr{
			Addr: IPv4Addr{127, 0, 0, 1},
			Mask: IPv4Addr{255, 255, 255, 255},
		}: "127.0.0.1/32",
		MaskedIPv4Addr{
			Addr: IPv4Addr{255, 255, 255, 255},
			Mask: IPv4Addr{255, 255, 128, 0},
		}: "255.255.255.255/17",
	}
	for a, s := range addrs {
		if a.String() != s {
			t.Errorf("invalid string for ipv4 address: actual=%v want=%v",
				a.String(), s)
		}
	}
}

func TestIPv6String(t *testing.T) {
	addrs := map[IPv6Addr]string{
		IPv6Addr{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}: "::1",
		IPv6Addr{0x20, 0x01, 0x0d, 0xb8, 0x12, 0x34}:             "2001:db8:1234::",
		IPv6Addr{0, 1, 0, 0, 0, 1}:                               "1:0:1::",
	}
	for a, s := range addrs {
		if a.String() != s {
			t.Errorf("invalid string for ipv6 address: actual=%v want=%v",
				a.String(), s)
		}
	}
}

func TestIPv6CIDR(t *testing.T) {
	addrs := map[IPv6Addr]int{
		IPv6Addr{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}: 0,
		IPv6Addr{255, 128}:                                       9,
		IPv6Addr{255, 255, 255, 255, 128}:                        33,
	}
	for a, m := range addrs {
		if a.AsCIDRMask() != m {
			t.Errorf("invalid CIDR for ipv6 mask: actual=%v want=%v",
				a.AsCIDRMask(), m)
		}
	}
}

func TestMaskedIPv6String(t *testing.T) {
	addrs := map[MaskedIPv6Addr]string{
		MaskedIPv6Addr{
			Addr: IPv6Addr{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			Mask: IPv6Addr{},
		}: "::1/0",
		MaskedIPv6Addr{
			Addr: IPv6Addr{0x02, 0x01, 0x0d, 0xb8, 0x12, 0x34},
			Mask: IPv6Addr{255, 255, 255, 255},
		}: "201:db8:1234::/32",
		MaskedIPv6Addr{
			Addr: IPv6Addr{0, 0, 255, 255, 255, 255},
			Mask: IPv6Addr{255, 255, 128, 0},
		}: "0:ffff:ffff::/17",
	}
	for a, s := range addrs {
		if a.String() != s {
			t.Errorf("invalid string for ipv6 address: actual=%v want=%v",
				a.String(), s)
		}
	}
}
