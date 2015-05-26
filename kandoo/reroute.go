package kandoo

import (
	"fmt"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
)

// ElephantDetected is a message emitted when an elephant flow is detected.
type ElephantDetected struct {
	Match nom.Match
}

// Rerouter implements the sample rerouting application in kandoo.
type Rerouter struct {
	NonLocal
}

func (r Rerouter) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	fmt.Printf("reroute: reroute an elephant flow %v\n",
		msg.Data().(ElephantDetected).Match)
	return nil
}
