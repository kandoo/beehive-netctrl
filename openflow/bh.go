package openflow

import (
	"flag"

	"github.com/kandoo/beehive/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/kandoo/beehive/bucket"

	bh "github.com/kandoo/beehive"
)

var (
	proto = flag.String("of.proto", "tcp", "protocol of the OpenFlow listener")
	addr  = flag.String("of.addr", "0.0.0.0:6633",
		"address of the OpenFlow listener in the form of HOST:PORT")
	readBufLen = flag.Int("of.rbuflen", 1<<8,
		"maximum number of packets to read per each read call")
	maxConnRate = flag.Int("of.maxrate", 1<<18,
		"maximum number of messages an openflow connection can generate per second")
)

// Option represents an OpenFlow listener option.
type Option func(l *ofListener)

// ListenOn returns an OpenFlow option that sets the address on which the
// OpenFlow service listens.
func ListenOn(addr string) Option {
	return func(l *ofListener) {
		l.addr = addr
	}
}

// UseProto returns an Openflow option that sets the protocol that the OpenFlow
// service uses to listen.
func UseProto(proto string) Option {
	return func(l *ofListener) {
		l.proto = proto
	}
}

// SetReadBufLen returns an OpenFlow option that sets reader buffer length of
// the OpenFlow service.
func SetReadBufLen(rlen int) Option {
	return func(l *ofListener) {
		l.readBufLen = rlen
	}
}

// StartOpenFlow starts the OpenFlow driver on the given hive using the default
// OpenFlow configuration that can be set through command line arguments.
func StartOpenFlow(hive bh.Hive, options ...Option) error {
	app := hive.NewApp("OFDriver",
		bh.OutRate(bucket.Rate(*maxConnRate), 10*uint64(*maxConnRate)))
	l := &ofListener{
		proto:      *proto,
		addr:       *addr,
		readBufLen: *readBufLen,
	}

	for _, opt := range options {
		opt(l)
	}

	app.Detached(l)
	glog.V(2).Infof("OpenFlow driver registered on %s:%s", l.proto, l.addr)

	return nil
}
