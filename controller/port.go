package controller

import (
	"fmt"

	"github.com/soheilhy/beehive-netctrl/nom"
	"github.com/soheilhy/beehive/bh"
)

type portStatusHandler struct{}

func (h *portStatusHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	// FIXME(soheil): This implementation is very naive and cannot tolerate
	// faults. We need to first check if the driver is the mater, and then apply
	// the change. Otherwise, we need to enque this message for that driver and
	// make sure we apply the log to the port status
	data := msg.Data().(nom.PortStatusChanged)
	dict := ctx.Dict(nodeDriversDict)
	k := bh.Key(data.Port.Node)
	n := nodeDrivers{}
	if err := dict.GetGob(k, &n); err != nil {
		return fmt.Errorf("Node %v not found", data.Port.Node)
	}

	if n.master() != data.Driver {
		return fmt.Errorf("%v is ignored since %v is not master", data.Port,
			data.Driver)
	}

	if p, ok := n.Ports.GetPort(data.Port.UID()); ok {
		if p == data.Port {
			return fmt.Errorf("Duplicate port status change for %v", data.Port)
		}

		n.Ports.DelPort(p)
	}

	n.Ports.AddPort(data.Port)
	ctx.Emit(nom.PortUpdated(data.Port))
	return nil
}

func (h *portStatusHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return bh.MappedCells{
		{nodeDriversDict, bh.Key(msg.Data().(nom.PortStatusChanged).Port.Node)},
	}
}
