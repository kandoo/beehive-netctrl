package openflow

import (
	"errors"

	"github.com/golang/glog"
	"github.com/soheilhy/beehive-netctrl/nom"
	"github.com/soheilhy/beehive-netctrl/openflow/of"
	"github.com/soheilhy/beehive-netctrl/openflow/of10"
	"github.com/soheilhy/beehive-netctrl/openflow/of12"
)

func (c *ofConn) handshake() (ofDriver, error) {
	hdr, err := c.ReadHeader()
	if err != nil {
		return nil, err
	}

	h, err := of.ToHello(hdr)
	if err != nil {
		return nil, err
	}

	glog.V(2).Info("Received hello from a switch")

	version := of.OPENFLOW_1_0
	if h.Version() > uint8(of.OPENFLOW_1_2) {
		version = of.OPENFLOW_1_2
	}
	h.SetVersion(uint8(version))

	if err = c.WriteHeader(h.Header); err != nil {
		return nil, err
	}
	c.Flush()

	glog.V(2).Info("Sent hello to the switch")

	var driver ofDriver
	switch version {
	case of.OPENFLOW_1_0:
		driver = &of10Driver{}
	case of.OPENFLOW_1_2:
		driver = &of12Driver{}
	}

	if err = driver.handshake(c); err != nil {
		return nil, err
	}

	if c.node.ID == nom.NodeID(0) {
		return nil, errors.New("Invalid node after handshake")
	}

	return driver, nil
}

func (d *of10Driver) handshake(c *ofConn) error {
	freq := of10.NewFeaturesRequest()
	if err := c.WriteHeader(freq.Header); err != nil {
		return err
	}
	c.Flush()

	glog.V(2).Info("Sent features request to the switch")

	hdr, err := c.ReadHeader()
	if err != nil {
		return err
	}

	v10, err := of10.ToHeader10(hdr)
	if err != nil {
		return err
	}

	frep, err := of10.ToFeaturesReply(v10)
	if err != nil {
		return err
	}

	glog.Infof("Handshake completed for switch %016x", frep.DatapathId())

	glog.Infof("Disabling packet buffers in the switch.")
	cfg := of10.NewSwitchSetConfig()
	cfg.SetMissSendLen(0xFFFF)
	c.WriteHeader(cfg.Header)

	nodeID := datapathIDToNodeID(frep.DatapathId())
	c.node = nom.Node{
		ID:           nodeID,
		Capabilities: nil,
	}

	c.ctx.Emit(nom.NodeConnected{
		Node: c.node,
		Driver: nom.Driver{
			BeeID: c.ctx.BeeID(),
			Role:  nom.DriverRoleDefault,
		},
	})

	for _, p := range frep.Ports() {
		glog.Infof("Port (switch=%016x, no=%d, mac=%012x, name=%s)\n",
			frep.DatapathId(), p.PortNo(), p.HwAddr(), p.Name())
		name := p.Name()
		c.ctx.Emit(nom.Port{
			ID:      portNoToPortID(uint32(p.PortNo())),
			Name:    string(name[:]),
			MACAddr: p.HwAddr(),
			Node:    c.NodeUID(),
		})
	}

	return nil
}

func (d *of12Driver) handshake(c *ofConn) error {
	freq := of12.NewFeaturesRequest()
	if err := c.WriteHeader(freq.Header); err != nil {
		return err
	}

	glog.V(2).Info("Sent features request to the switch")

	hdr, err := c.ReadHeader()
	if err != nil {
		return err
	}

	v12, err := of12.ToHeader12(hdr)
	if err != nil {
		return err
	}

	frep, err := of12.ToFeaturesReply(v12)
	if err != nil {
		return err
	}

	glog.Infof("Handshake completed for switch %016x", frep.DatapathId())

	glog.Infof("Disabling packet buffers in the switch.")
	cfg := of12.NewSwitchSetConfig()
	cfg.SetMissSendLen(0xFFFF)
	c.WriteHeader(cfg.Header)

	nodeID := datapathIDToNodeID(frep.DatapathId())
	c.node = nom.Node{
		ID: nodeID,
		Capabilities: []nom.NodeCapability{
			nom.CapDriverRole,
		},
	}

	c.ctx.Emit(nom.NodeConnected{
		Node: c.node,
		Driver: nom.Driver{
			BeeID: c.ctx.BeeID(),
			Role:  nom.DriverRoleDefault,
		},
	})

	for _, p := range frep.Ports() {
		glog.Infof("Port (switch=%016x, no=%d, mac=%012x, name=%s)\n",
			frep.DatapathId(), p.PortNo(), p.HwAddr(), p.Name())
		name := p.Name()
		c.ctx.Emit(nom.PortAdded{
			ID:      portNoToPortID(p.PortNo()),
			Name:    string(name[:]),
			MACAddr: p.HwAddr(),
			Node:    c.NodeUID(),
		})
	}

	return nil
}
