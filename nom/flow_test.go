package nom

import "testing"

func testMatch(t *testing.T, m1, m2 Match, subsumption, equality [2]bool) {
	if !m1.Equals(m1) {
		t.Errorf("%#v should be equal to itself", m1)
	}
	if !m2.Equals(m2) {
		t.Errorf("%#v should be equal to itself", m2)
	}
	if m1.Subsumes(m2) != subsumption[0] {
		t.Errorf("unexpected subsumption for %#v and %#v: actual=%v want=%v",
			m1, m2, !subsumption[0], subsumption[0])
	}
	if m2.Subsumes(m1) != subsumption[1] {
		t.Errorf("unexpected subsumption for %#v and %#v: actual=%v want=%v",
			m2, m1, !subsumption[1], subsumption[1])
	}
	if m1.Equals(m2) != equality[0] {
		t.Errorf("unexpected subsumption for %#v and %#v: actual=%v want=%v",
			m1, m2, !equality[0], equality[0])
	}
	if m2.Equals(m1) != equality[1] {
		t.Errorf("unexpected subsumption for %#v and %#v: actual=%v want=%v",
			m2, m1, !equality[1], equality[1])
	}
}

func TestInPort(t *testing.T) {
	m1 := Match{
		Fields: []Field{InPort(1)},
	}
	m2 := Match{
		Fields: []Field{InPort(2)},
	}
	testMatch(t, m1, m2, [2]bool{false, false}, [2]bool{false, false})
	testMatch(t, m1, m1, [2]bool{true, true}, [2]bool{true, true})
}

func TestEthSrc(t *testing.T) {
	m1 := Match{
		Fields: []Field{
			EthSrc{
				Addr: [6]byte{0x01, 0x02, 0x03, 0x00, 0x00, 0x00},
				Mask: [6]byte{0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00},
			},
		},
	}
	m2 := Match{
		Fields: []Field{
			EthSrc{
				Addr: [6]byte{0x01, 0x02, 0x03, 0x04, 0x00, 0x00},
				Mask: [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00},
			},
		},
	}
	testMatch(t, m1, m2, [2]bool{true, false}, [2]bool{false, false})
}

func TestEthDst(t *testing.T) {
	m1 := Match{
		Fields: []Field{
			EthDst{
				Addr: [6]byte{0x01, 0x02, 0x03, 0x00, 0x00, 0x00},
				Mask: [6]byte{0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00},
			},
		},
	}
	m2 := Match{
		Fields: []Field{
			EthDst{
				Addr: [6]byte{0x01, 0x02, 0x03, 0x04, 0x00, 0x00},
				Mask: [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00},
			},
		},
	}
	testMatch(t, m1, m2, [2]bool{true, false}, [2]bool{false, false})
}

func TestEth(t *testing.T) {
	m1 := Match{
		Fields: []Field{
			EthSrc{
				Addr: [6]byte{0x01, 0x02, 0x03, 0x00, 0x00, 0x00},
				Mask: [6]byte{0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00},
			},
			EthDst{
				Addr: [6]byte{0x01, 0x02, 0x03, 0x00, 0x00, 0x00},
				Mask: [6]byte{0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00},
			},
		},
	}
	m2 := Match{
		Fields: []Field{
			EthSrc{
				Addr: [6]byte{0x01, 0x02, 0x03, 0x00, 0x00, 0x00},
				Mask: [6]byte{0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00},
			},
			EthDst{
				Addr: [6]byte{0x01, 0x02, 0x03, 0x04, 0x00, 0x00},
				Mask: [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00},
			},
		},
	}
	testMatch(t, m1, m2, [2]bool{true, false}, [2]bool{false, false})
}

func TestIPv4Src(t *testing.T) {
	m1 := Match{
		Fields: []Field{
			IPv4Src{
				Addr: [4]byte{0x01, 0x02, 0x03, 0x00},
				Mask: [4]byte{0xFF, 0xFF, 0xFF, 0x00},
			},
		},
	}
	m2 := Match{
		Fields: []Field{
			IPv4Src{
				Addr: [4]byte{0x01, 0x02, 0x03, 0x04},
				Mask: [4]byte{0xFF, 0xFF, 0xFF, 0xFF},
			},
		},
	}
	testMatch(t, m1, m2, [2]bool{true, false}, [2]bool{false, false})
}

func TestIPv4Dst(t *testing.T) {
	m1 := Match{
		Fields: []Field{
			IPv4Dst{
				Addr: [4]byte{0x01, 0x02, 0x03, 0x00},
				Mask: [4]byte{0xFF, 0xFF, 0xFF, 0x00},
			},
		},
	}
	m2 := Match{
		Fields: []Field{
			IPv4Dst{
				Addr: [4]byte{0x01, 0x02, 0x03, 0x04},
				Mask: [4]byte{0xFF, 0xFF, 0xFF, 0xFF},
			},
		},
	}
	testMatch(t, m1, m2, [2]bool{true, false}, [2]bool{false, false})
}

func TestIPv4(t *testing.T) {
	m1 := Match{
		Fields: []Field{
			IPv4Src{
				Addr: [4]byte{0x01, 0x02, 0x03, 0x00},
				Mask: [4]byte{0xFF, 0xFF, 0xFF, 0x00},
			},
			IPv4Dst{
				Addr: [4]byte{0x01, 0x02, 0x03, 0x00},
				Mask: [4]byte{0xFF, 0xFF, 0xFF, 0x00},
			},
		},
	}
	m2 := Match{
		Fields: []Field{
			IPv4Src{
				Addr: [4]byte{0x01, 0x02, 0x03, 0x04},
				Mask: [4]byte{0xFF, 0xFF, 0xFF, 0xFF},
			},
			IPv4Dst{
				Addr: [4]byte{0x01, 0x02, 0x03, 0x00},
				Mask: [4]byte{0xFF, 0xFF, 0xFF, 0x00},
			},
		},
	}
	testMatch(t, m1, m2, [2]bool{true, false}, [2]bool{false, false})
}

func TestIPv6(t *testing.T) {
	m1 := Match{
		Fields: []Field{
			IPv6Src{
				Addr: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
				Mask: [16]byte{0xFF, 0xFF, 0xFF},
			},
			IPv6Dst{
				Addr: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
				Mask: [16]byte{0xFF, 0xFF, 0xFF, 0xFF},
			},
		},
	}
	m2 := Match{
		Fields: []Field{
			IPv6Src{
				Addr: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
				Mask: [16]byte{0xFF, 0xFF, 0xFF, 0xFF},
			},
			IPv6Dst{
				Addr: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
				Mask: [16]byte{0xFF, 0xFF, 0xFF, 0xFF},
			},
		},
	}
	testMatch(t, m1, m2, [2]bool{true, false}, [2]bool{false, false})
}

func TestTransportPortSrc(t *testing.T) {
	m1 := Match{
		Fields: []Field{TransportPortSrc(88)},
	}
	m2 := Match{
		Fields: []Field{TransportPortSrc(80)},
	}
	testMatch(t, m1, m2, [2]bool{false, false}, [2]bool{false, false})
	testMatch(t, m1, m1, [2]bool{true, true}, [2]bool{true, true})
}

func TestTransportPortDst(t *testing.T) {
	m1 := Match{
		Fields: []Field{TransportPortDst(88)},
	}
	m2 := Match{
		Fields: []Field{TransportPortDst(80)},
	}
	testMatch(t, m1, m2, [2]bool{false, false}, [2]bool{false, false})
	testMatch(t, m1, m1, [2]bool{true, true}, [2]bool{true, true})
}

func testFlowEntry(t *testing.T, f1, f2 FlowEntry,
	equality, subsumption [2]bool) {

	if f1.Equals(f2) != equality[0] {
		t.Errorf("invalid flow equality for %v and %v: actual=%v want=%v",
			f1, f2, !equality[0], equality[0])
	}
	if f2.Equals(f1) != equality[1] {
		t.Errorf("invalid flow equality for %v and %v: actual=%v want=%v",
			f2, f1, !equality[1], equality[1])
	}
	if f1.Subsumes(f2) != subsumption[0] {
		t.Errorf("invalid flow subsumption for %v and %v: actual=%v want=%v",
			f1, f2, !subsumption[0], subsumption[0])
	}
	if f2.Subsumes(f1) != subsumption[1] {
		t.Errorf("invalid flow subsumption for %v and %v: actual=%v want=%v",
			f2, f1, !subsumption[1], subsumption[1])
	}
}

func TestFlowEquality(t *testing.T) {
	f1 := FlowEntry{
		ID:   "0",
		Node: "n6",
		Match: Match{
			Fields: []Field{
				EthDst{
					Addr: MACAddr{0x1, 0x2, 0x3, 0x4, 0x5, 0x6},
					Mask: MACAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				},
				InPort("n6$$2"),
			},
		},
		Actions: []Action{
			ActionForward{Ports: []UID{"n6$$3"}},
		},
		Priority:    0x1,
		IdleTimeout: 0,
		HardTimeout: 0,
	}
	f2 := FlowEntry{
		ID:   "0",
		Node: "n6",
		Match: Match{
			Fields: []Field{
				EthDst{
					Addr: MACAddr{0x1, 0x2, 0x3, 0x4, 0x5, 0x6},
					Mask: MACAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				},
				InPort("n6$$2"),
			},
		},
		Actions: []Action{
			ActionForward{Ports: []UID{"n6$$3"}},
		},
		Priority:    0x1,
		IdleTimeout: 0,
		HardTimeout: 0,
	}
	testFlowEntry(t, f1, f1, [2]bool{true, true}, [2]bool{true, true})
	testFlowEntry(t, f1, f2, [2]bool{true, true}, [2]bool{true, true})

	f1.Match.Fields = f1.Match.Fields[1:]
	testFlowEntry(t, f1, f2, [2]bool{false, false}, [2]bool{true, false})

	f1.Match = f2.Match
	f2.Actions = append(f2.Actions, ActionFlood{})
	testFlowEntry(t, f1, f2, [2]bool{false, false}, [2]bool{false, false})
}
