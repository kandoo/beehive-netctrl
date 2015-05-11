package openflow

import (
	"io"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
	"github.com/kandoo/beehive-netctrl/openflow/of"
	"github.com/kandoo/beehive/Godeps/_workspace/src/github.com/golang/glog"
)

type ofConnConfig struct {
	readBufLen int
}

type ofConn struct {
	of.HeaderConn

	ctx bh.RcvContext // RcvContext of the detached bee running ofConn.

	readBufLen int         // Maximum number of packets to read.
	wCh        chan bh.Msg // Messages to be written.
	wErr       error       // Last error in write.

	node   nom.Node // Node that this connection represents.
	driver ofDriver // OpenFlow driver of this connection.
}

func (c *ofConn) drainWCh() {
	for {
		if _, ok := <-c.wCh; !ok {
			return
		}
	}
}

func (c *ofConn) Start(ctx bh.RcvContext) {
	defer func() {
		if c.driver != nil {
			c.driver.handleConnClose(c)
		}
		c.Close()
		// TODO(soheil): is there any better way to prevent deadlocks?
		glog.Infof("%v drains write queue for %v", ctx, c.RemoteAddr())
		go c.drainWCh()
	}()

	c.ctx = ctx
	c.wCh = make(chan bh.Msg, ctx.Hive().Config().DataChBufSize)

	var err error
	if c.driver, err = c.handshake(); err != nil {
		glog.Errorf("Error in OpenFlow handshake: %v", err)
		return
	}

	stop := make(chan struct{})

	wdone := make(chan struct{})
	go c.doWrite(wdone, stop)

	rdone := make(chan struct{})
	go c.doRead(rdone, stop)

	select {
	case <-rdone:
		close(stop)
	case <-wdone:
		close(stop)
	}

	<-rdone
	<-wdone
}

func (c *ofConn) doWrite(done chan struct{}, stop chan struct{}) {
	defer close(done)

	written := false
	var msg bh.Msg
	for {
		msg = nil
		if !written {
			select {
			case msg = <-c.wCh:
			case <-stop:
				return
			}
		} else {
			select {
			case msg = <-c.wCh:
			case <-stop:
				return
			default:
				if c.wErr = c.HeaderConn.Flush(); c.wErr != nil {
					return
				}
				written = false
				continue
			}
		}

		// Write the message.
		err := c.driver.handleMsg(msg, c)
		if c.wErr != nil {
			return
		}
		if err != nil {
			glog.Errorf("ofconn: Cannot convert NOM message to OpenFlow: %v",
				err)
		}
		written = true
	}
}

func (c *ofConn) doRead(done chan struct{}, stop chan struct{}) {
	defer close(done)

	pkts := make([]of.Header, c.readBufLen)
	for {
		select {
		case <-stop:
			return
		default:
		}

		n, err := c.ReadHeaders(pkts)
		if err != nil {
			if err == io.EOF {
				glog.Infof("connection %v closed", c.RemoteAddr())
			} else {
				glog.Errorf("cannot read from the connection %v: %v", c.RemoteAddr(),
					err)
			}
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
			pkts = make([]of.Header, c.readBufLen)
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
	c.wErr = c.HeaderConn.WriteHeader(pkt)
	return c.wErr
}
