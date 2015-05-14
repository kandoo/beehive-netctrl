package controller

import (
	"fmt"

	"github.com/kandoo/beehive/Godeps/_workspace/src/github.com/golang/glog"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
)

type portStatusHandler struct{}

func (h portStatusHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	// FIXME(soheil): This implementation is very naive and cannot tolerate
	// faults. We need to first check if the driver is the mater, and then apply
	// the change. Otherwise, we need to enque this message for that driver and
	// make sure we apply the log to the port status
	data := msg.Data().(nom.PortStatusChanged)
	dict := ctx.Dict(driversDict)
	k := string(data.Port.Node)
	n := nodeDrivers{}
	if err := dict.GetGob(k, &n); err != nil {
		return fmt.Errorf("NOMController: node %v not found", data.Port.Node)
	}

	if !n.isMaster(data.Driver) {
		glog.Warningf("NOMController: %v ignored, %v is not master, master is %v",
			data.Port, data.Driver, n.master())
		return nil
	}

	if p, ok := n.Ports.GetPort(data.Port.UID()); ok {
		if p == data.Port {
			return fmt.Errorf("NOMController: duplicate port status change for %v",
				data.Port)
		}

		n.Ports.DelPort(p)
	}

	n.Ports.AddPort(data.Port)
	ctx.Emit(nom.PortUpdated(data.Port))
	return nil
}

func (h portStatusHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return bh.MappedCells{
		{driversDict, string(msg.Data().(nom.PortStatusChanged).Port.Node)},
	}
}
