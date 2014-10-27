package openflow

import (
	"flag"

	"github.com/golang/glog"

	bh "github.com/kandoo/beehive"
)

// OFConfig stores the configuration of the OpenFlow driver.
type OFConfig struct {
	Proto      string // The driver's listening protocol.
	Addr       string // The driver's listening address.
	ReadBufLen int    // Maximum number of packets to read.
}

var defaultOFConfig = OFConfig{}

func init() {
	flag.StringVar(&defaultOFConfig.Proto, "ofproto", "tcp",
		"Protocol of the OpenFlow listener.")
	flag.StringVar(&defaultOFConfig.Addr, "ofaddr", "0.0.0.0:6633",
		"Address of the OpenFlow listener in the form of HOST:PORT.")
	flag.IntVar(&defaultOFConfig.ReadBufLen, "rbuflen", 1<<8,
		"Maximum number of packets to read per each read call.")
}

// StartOpenFlow starts the OpenFlow driver on the given hive using the default
// OpenFlow configuration that can be set through command line arguments.
func StartOpenFlow(hive bh.Hive) error {
	return StartOpenFlowWithConfig(hive, defaultOFConfig)
}

// StartOpenFlowWithConfig starts the OpenFlow driver on the give hive with the
// provided configuration.
func StartOpenFlowWithConfig(hive bh.Hive, cfg OFConfig) error {
	app := hive.NewApp("OFDriver")
	app.Detached(&ofListener{
		cfg: cfg,
	})

	glog.V(2).Infof("OpenFlow driver registered on %s:%s", cfg.Proto, cfg.Addr)
	return nil
}
