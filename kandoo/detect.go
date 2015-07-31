package kandoo

import (
	"fmt"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
)

type flow struct {
	Stats    nom.FlowStats
	Notified bool
}

// Detector implements the sample rerouting application in kandoo.
type Detector struct {
	Local

	ElephantThreshold uint64 // The minimum size of an elephent flow.
}

func (d Detector) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	res := msg.Data().(nom.FlowStatsQueryResult)

	var flows []flow
	if v, err := ctx.Dict("Switches").Get(string(res.Node)); err == nil {
		flows = v.([]flow)
	}

	for _, s1 := range res.Stats {
		// Find the previous stats.
		index := -1
		for i := range flows {
			if !s1.Match.Equals(flows[i].Stats.Match) {
				continue
			}

			if s1.Duration < flows[i].Stats.Duration {
				// The stat is newer than what we've seen so far.
				flows[i].Notified = false
			}

			flows[i].Stats = s1
			index = i
			break
		}

		if index < 0 {
			index = len(flows)
			flows = append(flows, flow{Stats: s1})
		}

		// Notify the router if it is an elephant flow, and make sure you do it
		// once.
		if !flows[index].Notified &&
			flows[index].Stats.Bytes > d.ElephantThreshold {

			flows[index].Notified = true
			ctx.Emit(ElephantDetected{
				Match: flows[index].Stats.Match,
			})
		}
	}

	return ctx.Dict("Switches").Put(string(res.Node), flows)
}

// Adder adds a node to the dictionary when it joins the network.
type Adder struct {
	Local
}

func (a Adder) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	joined := msg.Data().(nom.NodeJoined)

	return ctx.Dict("Switches").Put(string(nom.Node(joined).UID()), []flow{})
}

// Poller polls the switches.
type Poller struct {
	Local
}

func (p Poller) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	ctx.Dict("Switches").ForEach(func(k string, v interface{}) bool {
		fmt.Printf("poller: polling switch %v\n", k)
		ctx.Emit(nom.FlowStatsQuery{
			Node: nom.UID(k),
		})
		return true
	})
	return nil
}
