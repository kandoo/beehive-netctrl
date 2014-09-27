package openflow

import (
	"io"

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
	node   nom.Node
	driver ofDriver
	wCh    chan of.Header
}

func (c *ofConn) Start(ctx bh.RcvContext) {
	defer func() {
		c.driver.handleConnClose(c)
		c.Close()
	}()

	c.ctx = ctx
	c.wCh = make(chan of.Header, 4096)

	var err error
	if c.driver, err = c.handshake(); err != nil {
		glog.Errorf("Error in OpenFlow handshake: %v", err)
		return
	}

	pkts := make([]of.Header, c.cfg.readBufLen)
	for {
		n, err := c.ReadHeaders(pkts)
		if err != nil {
			if err == io.EOF {
				glog.Info("Connection closed.")
			} else {
				glog.Errorf("Cannot read from the connection: %v", err)
			}

			return
		}

		for _, pkt := range pkts[:n] {
			if err := c.driver.handlePkt(pkt, c); err != nil {
				glog.Errorf("%s", err)
				return
			}
		}

		var werr error
		for {
			select {
			case pkt := <-c.wCh:
				// If werr != nil, we loop and drain the wCh.
				if werr == nil {
					if werr = c.WriteHeader(pkt); werr != nil {
						glog.Errorf("ofconn: Cannot write packet: %v", werr)
					}
				}
			default:
				goto flush
			}
		}

	flush:
		if werr != nil {
			return
		}

		c.HeaderConn.Flush()

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
	c.wCh <- msg.Data().(of.Header)
	return nil
}

func (c *ofConn) NodeUID() nom.UID {
	return c.node.UID()
}
