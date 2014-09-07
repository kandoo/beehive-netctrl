package main

import (
	"github.com/soheilhy/beehive-netctrl/openflow"
	"github.com/soheilhy/beehive/bh"
)

func main() {
	h := bh.NewHive()
	openflow.StartOpenFlow(h)

	ch := make(chan bool)
	h.Start(ch)
}
