package discovery

import (
	"testing"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
)

func TestGraphBuilderCentralizedSinglePath(t *testing.T) {
	links := []nom.Link{
		{From: "n1$$1", To: "n2$$1"},
		{From: "n1$$2", To: "n3$$1"},
		{From: "n2$$2", To: "n4$$1"},
		{From: "n3$$2", To: "n5$$1"},
		{From: "n4$$2", To: "n5$$2"},
		{From: "n5$$3", To: "n6$$1"},
		{From: "n4$$3", To: "n6$$2"},
		{From: "n6$$1", To: "n5$$3"},
		{From: "n6$$2", To: "n4$$3"},
	}
	b := GraphBuilderCentralized{}
	ctx := &bh.MockRcvContext{}
	for _, l := range links {
		msg := &bh.MockMsg{
			MsgData: nom.LinkAdded(l),
		}
		b.Rcv(msg, ctx)
	}
	paths, l := ShortestPathCentralized("n1", "n6", ctx)
	if l != 3 {
		t.Errorf("invalid shortest path between n1 and n6: actual=%d want=3", l)
	}
	if len(paths) != 2 {
		t.Errorf("invalid number of paths between n1 and n6: actual=%d want=2",
			len(paths))
	}
	for _, p := range paths {
		if p[1] != links[2] && p[1] != links[3] {
			t.Errorf("invalid path: %v", p)
		}
	}
}
