package openflow

import (
	"testing"

	"github.com/kandoo/beehive-netctrl/nom"
)

func TestOF10Match(t *testing.T) {
	driver := of10Driver{}
	matches := []nom.Match{
		{
			Fields: []nom.Field{
				nom.EthSrc{
					Addr: nom.MACAddr{1, 2, 3, 4, 5, 6},
					Mask: nom.MaskNoneMAC,
				},
			},
		},
		{
			Fields: []nom.Field{
				nom.EthDst{
					Addr: nom.MACAddr{1, 2, 3, 4, 5, 6},
					Mask: nom.MaskNoneMAC,
				},
			},
		},
		{
			Fields: []nom.Field{
				nom.IPv4Src{
					Addr: nom.IPv4Addr{1, 2, 3, 4},
					Mask: nom.IPv4Addr{255, 255, 255, 0},
				},
			},
		},
		{
			Fields: []nom.Field{
				nom.IPv4Dst{
					Addr: nom.IPv4Addr{127, 0, 0, 1},
					Mask: nom.IPv4Addr{255, 255, 255, 128},
				},
			},
		},
	}
	for _, m := range matches {
		ofm, err := driver.ofMatch(m)
		if err != nil {
			t.Error(err)
		}
		nm, err := driver.nomMatch(ofm)
		if err != nil {
			t.Error(err)
		}
		if !nm.Equals(m) {
			t.Errorf("invalid match conversion:\n\tactual=%#v\n\twant=%#v", nm, m)
		}
	}
}
