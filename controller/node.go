package controller

import (
	"fmt"

	"github.com/kandoo/beehive/Godeps/_workspace/src/github.com/golang/glog"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
)

type nodeConnectedHandler struct{}

func (h nodeConnectedHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	nc := msg.Data().(nom.NodeConnected)

	dict := ctx.Dict(driversDict)
	k := string(nc.Node.ID)
	n := nodeDrivers{}
	if err := nom.DictGet(dict, k, &n); err != nil {
		n.Node = nc.Node
		n.Drivers = nom.Drivers{nc.Driver}
		if err := dict.PutGob(k, &n); err != nil {
			return err
		}

		glog.V(2).Infof("%v joins", nc.Node)
		ctx.Emit(nom.NodeJoined(nc.Node))
		return nil
	}

	if n.hasDriver(nc.Driver) {
		return fmt.Errorf("driver %v reconnects to %v", nc.Driver, n.Node)
	}

	n.Drivers = append(n.Drivers, nc.Driver)
	return dict.PutGob(k, &n)
}

func (h nodeConnectedHandler) Map(msg bh.Msg,
	ctx bh.MapContext) bh.MappedCells {

	nc := msg.Data().(nom.NodeConnected)
	return bh.MappedCells{{driversDict, string(nc.Node.ID)}}
}

type nodeDisconnectedHandler struct{}

func (h nodeDisconnectedHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	dc := msg.Data().(nom.NodeDisconnected)
	d := ctx.Dict(driversDict)
	k := string(dc.Node.ID)
	nd := nodeDrivers{}
	if err := d.GetGob(k, &nd); err != nil {
		return fmt.Errorf("driver %v disconnects from %v before connecting",
			dc.Driver, dc.Node)
	}

	if !nd.removeDriver(dc.Driver) {
		return fmt.Errorf("driver %v disconnects from %v before connecting",
			dc.Driver, dc.Node)
	}

	if len(nd.Drivers) == 0 {
		ctx.Emit(nom.NodeLeft(nd.Node))
		return d.Del(k)
	}
	return d.PutGob(k, &nd)
}

func (h nodeDisconnectedHandler) Map(msg bh.Msg,
	ctx bh.MapContext) bh.MappedCells {

	nd := msg.Data().(nom.NodeDisconnected)
	return bh.MappedCells{{driversDict, string(nd.Node.ID)}}
}
