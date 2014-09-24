package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/soheilhy/beehive-netctrl/openflow"
	"github.com/soheilhy/beehive/bh"
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

	h.Start()
}
