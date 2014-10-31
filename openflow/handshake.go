package openflow

import (
	"errors"

	"github.com/kandoo/beehive/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/kandoo/beehive-netctrl/nom"
	"github.com/kandoo/beehive-netctrl/openflow/of"
	"github.com/kandoo/beehive-netctrl/openflow/of10"
	"github.com/kandoo/beehive-netctrl/openflow/of12"
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
		MACAddr:      datapathIDToMACAddr(frep.DatapathId()),
		Capabilities: nil,
	}
	glog.Infof("%v connected", c.node)

	nomDriver := nom.Driver{
		BeeID: c.ctx.ID(),
		Role:  nom.DriverRoleDefault,
	}

	c.ctx.Emit(nom.NodeConnected{
		Node:   c.node,
		Driver: nomDriver,
	})

	d.ofPorts = make(map[uint16]*nom.Port)
	d.nomPorts = make(map[nom.UID]uint16)
	for _, p := range frep.Ports() {
		name := p.Name()
		port := nom.Port{
			ID:      portNoToPortID(uint32(p.PortNo())),
			Name:    string(name[:]),
			MACAddr: p.HwAddr(),
			Node:    c.NodeUID(),
		}
		d.ofPorts[p.PortNo()] = &port
		d.nomPorts[port.UID()] = p.PortNo()
		glog.Infof("%v added", port)
		if p.PortNo() <= uint16(of10.PP_MAX) {
			c.ctx.Emit(nom.PortStatusChanged{
				Port:   port,
				Driver: nomDriver,
			})
		}
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
		ID:      nodeID,
		MACAddr: datapathIDToMACAddr(frep.DatapathId()),
		Capabilities: []nom.NodeCapability{
			nom.CapDriverRole,
		},
	}

	nomDriver := nom.Driver{
		BeeID: c.ctx.ID(),
		Role:  nom.DriverRoleDefault,
	}
	c.ctx.Emit(nom.NodeConnected{
		Node:   c.node,
		Driver: nomDriver,
	})

	d.ofPorts = make(map[uint32]*nom.Port)
	d.nomPorts = make(map[nom.UID]uint32)
	for _, p := range frep.Ports() {
		if p.PortNo() > uint32(of12.PP_MAX) {
			continue
		}
		name := p.Name()
		port := nom.Port{
			ID:      portNoToPortID(p.PortNo()),
			Name:    string(name[:]),
			MACAddr: p.HwAddr(),
			Node:    c.NodeUID(),
		}
		d.ofPorts[p.PortNo()] = &port
		d.nomPorts[port.UID()] = p.PortNo()
		glog.Infof("%v added", port)
		c.ctx.Emit(nom.PortStatusChanged{
			Port:   port,
			Driver: nomDriver,
		})
	}

	return nil
}
