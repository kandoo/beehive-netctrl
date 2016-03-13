package main

import (
	"flag"
	"os"
	"runtime/pprof"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/controller"
	"github.com/kandoo/beehive-netctrl/discovery"
	"github.com/kandoo/beehive-netctrl/openflow"
	"github.com/kandoo/beehive/Godeps/_workspace/src/github.com/golang/glog"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			glog.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	h := bh.NewHive()
	openflow.StartOpenFlow(h)
	controller.RegisterNOMController(h)
	discovery.RegisterDiscovery(h)

	// Register a switch:
	// switching.RegisterSwitch(h, bh.Persistent(1))
	// or a hub:
	// switching.RegisterHub(h, bh.NonTransactional())

	h.Start()
}
