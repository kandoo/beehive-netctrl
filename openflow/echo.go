package openflow

import (
	"github.com/golang/glog"
	"github.com/soheilhy/beehive-netctrl/openflow/of"
	"github.com/soheilhy/beehive-netctrl/openflow/of10"
	"github.com/soheilhy/beehive-netctrl/openflow/of12"
)

func (d *of10Driver) handleEchoRequest(req of10.EchoRequest, c *ofConn) error {
	return doHandleEchoRequest(req.Header, of10.NewEchoReply().Header, c)
}

func (d *of12Driver) handleEchoRequest(req of12.EchoRequest, c *ofConn) error {
	return doHandleEchoRequest(req.Header, of12.NewEchoReply().Header, c)
}

func doHandleEchoRequest(req of.Header, res of.Header, c *ofConn) error {
	glog.V(2).Infof("Received an echo request from the switch")
	res.SetXid(req.Xid())
	err := c.WriteHeaders([]of.Header{res})
	if err != nil {
		return err
	}
	glog.V(2).Infof("Sent an echo reply from the switch")
	return nil
}
