package openflow

import (
	"errors"
	"net"

	"github.com/golang/glog"
	"github.com/soheilhy/beehive-netctrl/openflow/of"
	"github.com/soheilhy/beehive/bh"
)

type ofListener struct {
	cfg OFConfig
}

func (l *ofListener) Start(ctx bh.RcvContext) {
	nl, err := net.Listen(l.cfg.Proto, l.cfg.Addr)
	if err != nil {
		glog.Errorf("Cannot start the OF listener: %v", err)
		return
	}

	glog.Infof("OF listener started on %s:%s", l.cfg.Proto, l.cfg.Addr)

	defer func() {
		glog.Infof("OF listener closed")
		nl.Close()
	}()

	for {
		c, err := nl.Accept()
		if err != nil {
			glog.Errorf("Error in OF accept: %v", err)
			return
		}

		l.startOFConn(c, ctx)
	}
}

func (l *ofListener) startOFConn(conn net.Conn, ctx bh.RcvContext) {
	ofc := &ofConn{
		HeaderConn: of.NewHeaderConn(conn),
		cfg: ofConnConfig{
			readBufLen: l.cfg.ReadBufLen,
		},
	}

	ctx.StartDetached(ofc)
}

func (l *ofListener) Stop(ctx bh.RcvContext) {
}

func (l *ofListener) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	return errors.New("No message should be sent to the listener")
}
