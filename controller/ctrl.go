package controller

import (
	"github.com/soheilhy/beehive-netctrl/nom"
	"github.com/soheilhy/beehive/bh"
)

func RegisterNOMController(h bh.Hive) {
	app := h.NewApp("NOMController")

	app.Handle(nom.NodeConnected{}, &nodeConnectedHandler{})
	app.Handle(nom.NodeDisconnected{}, &nodeDisconnectedHandler{})

	app.Handle(nom.AddFlowEntry{}, &addFlowHandler{})
	app.Handle(nom.DelFlowEntry{}, &delFlowHandler{})

	app.Handle(nom.FlowStatQuery{}, &queryHandler{})

	app.Handle(nom.PacketOut{}, &pktOutHandler{})
}
