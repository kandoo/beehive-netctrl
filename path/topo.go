package path

import (
	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/discovery"
	"github.com/kandoo/beehive-netctrl/nom"
	"github.com/kandoo/beehive/Godeps/_workspace/src/github.com/golang/glog"
)

func inPortsFromOutPorts(outport []nom.UID, ctx bh.RcvContext) (
	inports []nom.UID) {

nextoutport:
	for _, p := range outport {
		n := nom.NodeFromPortUID(p)
		for _, l := range discovery.LinksCentralized(n, ctx) {
			if l.From == p {
				inports = append(inports, p)
				continue nextoutport
			}
		}
		glog.Errorf("cannot find peer port for %v", p)
	}
	return inports
}

func outPortsFromFloodNode(n, inp nom.UID, ctx bh.RcvContext) (
	outports []nom.UID) {

	for _, l := range discovery.LinksCentralized(n, ctx) {
		if l.From != inp {
			continue
		}
		outports = append(outports, l.From)
	}
	return outports
}
