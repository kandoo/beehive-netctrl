package controller

import (
	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
)

func RegisterNOMController(h bh.Hive) {
	app := h.NewApp("NOMController")

	app.Handle(nom.NodeConnected{}, nodeConnectedHandler{})
	app.Handle(nom.NodeDisconnected{}, nodeDisconnectedHandler{})
	app.Handle(nom.PortStatusChanged{}, portStatusHandler{})

	app.Handle(nom.AddFlowEntry{}, addFlowHandler{})
	app.Handle(nom.DelFlowEntry{}, delFlowHandler{})

	app.Handle(nom.FlowStatQuery{}, queryHandler{})

	app.Handle(nom.PacketOut{}, pktOutHandler{})
}
