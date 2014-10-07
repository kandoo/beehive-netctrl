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
	cfg    ofConnConfig  // Configuration of this connection.
	ctx    bh.RcvContext // RcvContext of the detached bee running the connection.
	node   nom.Node      // Node that this connection represents.
	driver ofDriver      // OpenFlow driver of this connection.
	wCh    chan bh.Msg   // Messages to be written.
	wErr   error         // Last error in write.
}

func (c *ofConn) Start(ctx bh.RcvContext) {
	defer func() {
		if c.driver != nil {
			c.driver.handleConnClose(c)
		}
		c.Close()
	}()

	c.ctx = ctx
	c.wCh = make(chan bh.Msg, 4096)

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

		for {
			select {
			case msg := <-c.wCh:
				if c.wErr != nil {
					// Drain the channel.
					continue
				}

				// If werr != nil, we loop and drain the wCh.
				if err := c.driver.handleMsg(msg, c); err != nil {
					glog.Errorf("ofconn: Cannot convert NOM message to OpenFlow: %v",
						err)
					continue
				}
			default:
				goto flush
			}
		}

	flush:
		if c.wErr != nil {
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
	c.wCh <- msg
	return nil
}

func (c *ofConn) NodeUID() nom.UID {
	return c.node.UID()
}

func (c *ofConn) WriteHeader(pkt of.Header) error {
	if c.wErr != nil {
		return c.wErr
	}

	c.wErr = c.HeaderConn.WriteHeader(pkt)
	return c.wErr
}
