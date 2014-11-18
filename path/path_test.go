package path

import (
	"testing"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/discovery"
	"github.com/kandoo/beehive-netctrl/nom"
)

func buildTopologyForTest() *bh.MockRcvContext {
	links := []nom.Link{
		{From: "n1$$1", To: "n2$$1"},
		{From: "n2$$1", To: "n1$$1"},
		{From: "n2$$2", To: "n4$$1"},
		{From: "n4$$1", To: "n2$$2"},
		{From: "n3$$2", To: "n5$$1"},
		{From: "n5$$1", To: "n3$$2"},
		{From: "n4$$2", To: "n6$$1"},
		{From: "n6$$1", To: "n4$$2"},
		{From: "n5$$2", To: "n6$$2"},
		{From: "n6$$2", To: "n5$$2"},
	}
	b := discovery.GraphBuilderCentralized{}
	ctx := &bh.MockRcvContext{}
	for _, l := range links {
		msg := &bh.MockMsg{
			MsgData: nom.LinkAdded(l),
		}
		b.Rcv(msg, ctx)
	}
	return ctx
}

func TestAddP2PPath(t *testing.T) {
	ctx := buildTopologyForTest()
	p := addHandler{}
	msg := &bh.MockMsg{
		MsgData: nom.AddPath{},
	}
	err := p.Rcv(msg, ctx)
	if err == nil {
		t.Error("no error on invalid path")
	}
	msg.MsgData = nom.AddPath{
		Path: nom.Path{
			Pathlets: []nom.Pathlet{
				{
					Match: nom.Match{
						Fields: []nom.Field{
							nom.InPort("n1$$0"),
						},
					},
					Actions: []nom.Action{
						nom.ActionForward{
							Ports: []nom.UID{"n6$$3"},
						},
					},
				},
			},
			Priority: 1,
		},
	}
	if err := p.Rcv(msg, ctx); err != nil {
		t.Errorf("cannot install flows for path: %v", err)
	}
	if len(ctx.CtxMsgs) == 0 {
		t.Error("no flows installed")
	}

	iports := []nom.UID{"n1$$0", "n2$$1", "n4$$1", "n6$$1"}
	oports := []nom.UID{"n1$$1", "n2$$2", "n4$$2", "n6$$3"}
	for i, msg := range ctx.CtxMsgs {
		add := msg.Data().(nom.AddFlowEntry)
		if add.Flow.Priority != 1 {
			t.Errorf("invalid flow priority: actual=%v want=1", add.Flow.Priority)
		}
		iport, ok := add.Flow.Match.InPort()
		if !ok {
			t.Errorf("flow #%v has no in ports", i)
		} else if nom.UID(iport) != iports[i] {
			t.Errorf("invalid input port on flow #%v: actual=%v want=%v", i, iport,
				iports[i])
		}

		oport := add.Flow.Actions[0].(nom.ActionForward).Ports[0]
		if oport != oports[i] {
			t.Errorf("invalid output port on flow #%v: actual=%v want=%v", i, oport,
				oports[i])
		}
	}
}

func TestAddL2Path(t *testing.T) {
	ctx := buildTopologyForTest()
	p := addHandler{}
	msg := &bh.MockMsg{
		MsgData: nom.AddPath{},
	}
	err := p.Rcv(msg, ctx)
	if err == nil {
		t.Error("no error on invalid path")
	}
	ethDst := nom.EthDst{
		Addr: nom.MACAddr{1, 2, 3, 4, 5, 6},
		Mask: nom.MaskNoneMAC,
	}
	msg.MsgData = nom.AddPath{
		Path: nom.Path{
			Pathlets: []nom.Pathlet{
				{
					Match: nom.Match{
						Fields: []nom.Field{
							ethDst,
						},
					},
					Actions: []nom.Action{
						nom.ActionForward{
							Ports: []nom.UID{"n6$$3"},
						},
					},
				},
			},
			Priority: 1,
		},
	}
	if err := p.Rcv(msg, ctx); err != nil {
		t.Errorf("cannot install flows for path: %v", err)
	}
	if len(ctx.CtxMsgs) == 0 {
		t.Error("no flows installed")
	}

	out := map[nom.UID]nom.UID{
		"n1": "n1$$1",
		"n2": "n2$$2",
		"n3": "n3$$2",
		"n4": "n4$$2",
		"n5": "n5$$2",
		"n6": "n6$$3",
	}
	for i, msg := range ctx.CtxMsgs {
		add := msg.Data().(nom.AddFlowEntry)
		if add.Flow.Priority != 1 {
			t.Errorf("invalid flow priority: actual=%v want=1", add.Flow.Priority)
		}
		fethDst, ok := add.Flow.Match.EthDst()
		if !ok {
			t.Errorf("flow #%v has no destination eth addr", i)
		} else if !fethDst.Equals(ethDst) {
			t.Errorf("invalid eth address for flow #%v: actual=%v want=%v", i,
				fethDst, ethDst)
		}

		delete(out, add.Flow.Node)
	}

	for _, p := range out {
		t.Errorf("no flow installed on port %v", p)
	}
}
