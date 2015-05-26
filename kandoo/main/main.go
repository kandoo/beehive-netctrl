package main

import (
	"flag"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/controller"
	"github.com/kandoo/beehive-netctrl/discovery"
	"github.com/kandoo/beehive-netctrl/kandoo"
	"github.com/kandoo/beehive-netctrl/openflow"
	"github.com/kandoo/beehive-netctrl/switching"
)

var eThreshold = flag.Uint64("kandoo.thresh", 1024,
	"the minimum size of an elephent flow ")

func main() {
	h := bh.NewHive()

	openflow.StartOpenFlow(h)
	controller.RegisterNOMController(h)
	discovery.RegisterDiscovery(h)
	switching.RegisterSwitch(h)
	kandoo.RegisterApps(h, *eThreshold)

	h.Start()
}
