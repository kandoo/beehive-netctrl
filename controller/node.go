package controller

import (
	"fmt"
	"time"

	"github.com/kandoo/beehive/Godeps/_workspace/src/github.com/golang/glog"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
)

type nodeConnectedHandler struct{}

func (h nodeConnectedHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	nc := msg.Data().(nom.NodeConnected)

	ddict := ctx.Dict(driversDict)
	k := string(nc.Node.ID)
	n := nodeDrivers{}
	if err := nom.DictGet(ddict, k, &n); err != nil {
		n.Node = nc.Node
	}

	if _, ok := n.driver(nc.Driver); ok {
		return fmt.Errorf("driver %v reconnects to %v", nc.Driver, n.Node)
	}

	gdict := ctx.Dict(genDict)
	gen := uint64(0)
	gdict.GetGob("gen", &gen)
	gen++

	db := nc.Driver.BeeID
	if len(n.Drivers) == 0 {
		nc.Driver.Role = nom.DriverRoleMaster
		ctx.Emit(nom.NodeJoined(nc.Node))
		glog.V(2).Infof("%v connected to master controller", nc.Node)
	} else {
		nc.Driver.Role = nom.DriverRoleSlave
		glog.V(2).Infof("%v connected to slave controller", nc.Node)
	}
	n.Drivers = append(n.Drivers, driverInfo{
		Driver:   nc.Driver,
		LastSeen: time.Now(),
	})

	ctx.SendToBee(nom.ChangeDriverRole{
		Node:       nc.Node.UID(),
		Role:       nc.Driver.Role,
		Generation: gen,
	}, db)

	gdict.PutGob("gen", gen)
	return ddict.PutGob(k, &n)
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

	ndd, ok := nd.driver(dc.Driver)
	if !ok {
		return fmt.Errorf("driver %v disconnects from %v before connecting",
			dc.Driver, dc.Node)
	}
	nd.removeDriver(dc.Driver)

	if len(nd.Drivers) == 0 {
		ctx.Emit(nom.NodeLeft(nd.Node))
		return d.Del(k)
	}

	if ndd.Role == nom.DriverRoleMaster {
		// TODO(soheil): Maybe a smarter load balacning technique.
		gdict := ctx.Dict(genDict)
		gen := uint64(0)
		gdict.GetGob("gen", &gen)
		gen++

		nd.Drivers[0].Role = nom.DriverRoleMaster
		ctx.SendToBee(nom.ChangeDriverRole{
			Node:       nd.Node.UID(),
			Role:       nom.DriverRoleMaster,
			Generation: gen,
		}, nd.master().BeeID)
	}

	return d.PutGob(k, &nd)
}

func (h nodeDisconnectedHandler) Map(msg bh.Msg,
	ctx bh.MapContext) bh.MappedCells {

	nd := msg.Data().(nom.NodeDisconnected)
	return bh.MappedCells{{driversDict, string(nd.Node.ID)}}
}

type roleUpdateHandler struct{}

func (h roleUpdateHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	ru := msg.Data().(nom.DriverRoleUpdate)

	ddict := ctx.Dict(driversDict)
	k := string(ru.Node)
	n := nodeDrivers{}
	if err := nom.DictGet(ddict, k, &n); err != nil {
		return fmt.Errorf("role update received before node %v connects", ru.Node)
	}

	found := false
	for i := range n.Drivers {
		if n.Drivers[i].BeeID == ru.Driver.BeeID {
			found = true
			n.Drivers[i].Role = ru.Driver.Role
			break
		}
	}

	if !found {
		return fmt.Errorf("role update received before driver %v connects",
			ru.Driver)
	}

	gdict := ctx.Dict(genDict)
	gen := uint64(0)
	gdict.GetGob("gen", &gen)
	if ru.Generation > gen {
		gdict.PutGob("gen", ru.Generation)
	}

	return nil
}

func (h roleUpdateHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	ru := msg.Data().(nom.DriverRoleUpdate)
	return bh.MappedCells{{driversDict, string(ru.Node)}}
}
