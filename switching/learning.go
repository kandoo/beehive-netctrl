package switching

import (
	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
	"github.com/kandoo/beehive/Godeps/_workspace/src/github.com/golang/glog"
)

type LearningSwitch struct {
	Hub
}

func (h LearningSwitch) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	in := msg.Data().(nom.PacketIn)
	src := in.Packet.SrcMAC()
	dst := in.Packet.DstMAC()
	glog.V(2).Infof("received packet in from %v to %v", src, dst)
	if dst.IsMulticast() {
		// TODO(soheil): just drop LLDP.
		glog.Infof("dropped multi-cast packet to %v", dst)
		return nil
	}

	if dst.IsBroadcast() {
		return h.Hub.Rcv(msg, ctx)
	}

	d := ctx.Dict("mac2port")
	srck := src.Key()
	var p nom.UID
	if err := d.GetGob(srck, &p); err != nil || p != in.InPort {
		if err == nil {
			// TODO(soheil): maybe add support for multi ports.
			glog.Infof("%v is moved from port %v to port %v", src, p, in.InPort)
		}

		if err = d.PutGob(srck, &in.InPort); err != nil {
			glog.Fatalf("cannot serialize port: %v", err)
		}
	}

	dstk := dst.Key()
	err := d.GetGob(dstk, &p)
	if err != nil {
		return h.Hub.Rcv(msg, ctx)
	}

	add := nom.AddFlowEntry{
		Flow: nom.FlowEntry{
			Node: in.Node,
			Match: nom.Match{
				Fields: []nom.Field{
					nom.EthDst{Addr: dst},
				},
			},
			Actions: []nom.Action{
				nom.ActionForward{
					Ports: []nom.UID{p},
				},
			},
		},
	}
	ctx.ReplyTo(msg, add)

	out := nom.PacketOut{
		Node:     in.Node,
		InPort:   in.InPort,
		BufferID: in.BufferID,
		Packet:   in.Packet,
		Actions: []nom.Action{
			nom.ActionForward{
				Ports: []nom.UID{p},
			},
		},
	}
	ctx.ReplyTo(msg, out)
	return nil
}
