package kandoo

import (
	"time"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
)

// RegisterApps registers Kandoo applications on the hive, with the given
// elephant flow size threshold.
func RegisterApps(hive bh.Hive, threshold uint64) {
	ar := hive.NewApp("Reroute")
	ar.Handle(ElephantDetected{}, Rerouter{})

	ad := hive.NewApp("Detect")
	ad.Handle(nom.FlowStatsQueryResult{}, Detector{})
	ad.Handle(nom.NodeJoined{}, Adder{})

	type poll struct{}
	ad.Handle(poll{}, Poller{})
	ad.Detached(bh.NewTimer(1*time.Second, func() {
		hive.Emit(poll{})
	}))
}
