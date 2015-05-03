package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/controller"
	"github.com/kandoo/beehive-netctrl/discovery"
	"github.com/kandoo/beehive-netctrl/openflow"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	h := bh.NewHive()
	openflow.StartOpenFlow(h)
	controller.RegisterNOMController(h)
	discovery.RegisterDiscovery(h)

	// app := h.NewApp("Hub", bh.NonTransactional())
	// app.Handle(nom.PacketIn{}, &switching.Hub{})

	h.Start()
}
