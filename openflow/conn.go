package openflow

import (
	"github.com/golang/glog"
	"github.com/soheilhy/beehive-netctrl/nom"
	"github.com/soheilhy/beehive-netctrl/openflow/of"
	"github.com/soheilhy/beehive/bh"
)

type ofConnConfig struct {
	readBufLen int
}

type ofConn struct {
	of.HeaderConn
	cfg    ofConnConfig
	ctx    bh.RcvContext
	node   nom.NodeID
	driver ofDriver
}

func (c *ofConn) Start(ctx bh.RcvContext) {
	defer func() {
		c.Close()
	}()

	var err error
	if c.driver, err = c.handshake(); err != nil {
		glog.Errorf("Error in OpenFlow handshake: %v", err)
		return
	}

	pkts := make([]of.Header, c.cfg.readBufLen)
	for {
		n, err := c.ReadHeaders(pkts)
		if err != nil {
			glog.Errorf("Cannot read from the connection: %v", err)
			return
		}

		for _, pkt := range pkts[:n] {
			if err := c.driver.handlePkt(pkt, c); err != nil {
				glog.Errorf("%s", err)
				return
			}
		}

		pkts = pkts[n:]
		if len(pkts) == 0 {
			pkts = make([]of.Header, c.cfg.readBufLen)
		}
	}
}

func (c *ofConn) Stop(ctx bh.RcvContext) {
	c.Close()
}

func (c *ofConn) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	pkt := msg.Data().(of.Header)
	return c.WriteHeader(pkt)
}
