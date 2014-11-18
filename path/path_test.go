package path

import (
	"testing"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/discovery"
	"github.com/kandoo/beehive-netctrl/nom"
)

func TestAddPath(t *testing.T) {
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
