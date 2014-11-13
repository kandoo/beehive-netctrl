package openflow

import (
	"fmt"
	"time"

	"github.com/kandoo/beehive-netctrl/nom"
	"github.com/kandoo/beehive-netctrl/openflow/of10"
	"github.com/kandoo/beehive/Godeps/_workspace/src/github.com/golang/glog"
)

func (d *of10Driver) handleStatsReply(reply of10.StatsReply,
	c *ofConn) error {

	switch {
	case of10.IsFlowStatsReply(reply):
		return d.handleFlowStatsReply(of10.NewFlowStatsReplyWithBuf(reply.Buf), c)
	default:
		return fmt.Errorf("of10Driver: unsupported stats type %v",
			reply.StatsType())
	}
}

func (d *of10Driver) handleFlowStatsReply(reply of10.FlowStatsReply,
	c *ofConn) error {

	nomReply := nom.FlowStatsQueryResult{
		Node: c.node.UID(),
	}
	for _, stat := range reply.FlowStats() {
		m, err := d.nomMatch(stat.Match())
		if err != nil {
			return err
		}
		stat := nom.FlowStats{
			Match: m,
			Duration: time.Duration(stat.DurationSec())*time.Second +
				time.Duration(stat.DurationNsec()),
			Packets: stat.PacketCount(),
			Bytes:   stat.ByteCount(),
		}
		nomReply.Stats = append(nomReply.Stats, stat)
	}
	c.ctx.Emit(nomReply)
	return nil
}

func (d *of12Driver) handleFlowStatsReply(reply of10.FlowStatsReply,
	c *ofConn) error {

	glog.Fatalf("TODO stat reply not implemented")
	return nil
}
