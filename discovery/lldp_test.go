package discovery

import (
	"testing"

	"github.com/kandoo/beehive-netctrl/nom"
)

func TestLLDPEncode(t *testing.T) {
	node := nom.Node{
		ID:      "n1",
		MACAddr: [6]byte{1, 2, 3, 4, 5, 6},
	}
	port := nom.Port{
		ID:      "1",
		Name:    "if1",
		MACAddr: [6]byte{1, 2, 3, 4, 5, 7},
		Node:    node.UID(),
	}
	b := encodeLLDP(node, port)
	//for i := 0; i < len(b); i += 12 {
	//j := 0
	//for ; j < 12 && i+j < len(b); j++ {
	//fmt.Printf("%02X ", b[i+j])
	//}
	//for k := j; k < 12; k++ {
	//fmt.Print("   ")
	//}
	//fmt.Printf("\t%s\n", strconv.QuoteToASCII(string(b[i:i+j])))
	//}

	decN, decP, err := decodeLLDP(b)
	if err != nil {
		t.Errorf("Cannot decode LLDP: %v", err)
	}

	if decN.ID != node.ID {
		t.Errorf("Invalid node ID decoded: %v != %v", decN.ID, node.ID)
	}

	if decN.MACAddr != node.MACAddr {
		t.Errorf("Invalid node MAC decoded: %v != %v", decN.MACAddr, node.MACAddr)
	}

	if decP.ID != port.ID {
		t.Errorf("Invalid port ID decoded: %v != %v", decP.ID, port.ID)
	}
}
