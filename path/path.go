package path

import (
	"errors"
	"fmt"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/discovery"
	"github.com/kandoo/beehive-netctrl/nom"
)

type addHandler struct{}

func (h addHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	add := msg.Data().(nom.AddPath)

	if len(add.Path.Pathlets) == 0 {
		return errors.New("path: path has no pathlets")
	}

	flows := make([]nom.FlowEntry, 0, len(add.Path.Pathlets))
	// TODO(soheil): maybe detect loops in pathlets?
	var outports []nom.UID
	var newflows []nom.FlowEntry
	var err error
	for _, p := range add.Path.Pathlets {
		if len(outports) == 0 {
			newflows, outports, err = genFlowsForPathlet(p, nom.Nil,
				add.Path.Priority, ctx)
			if err != nil {
				return err
			}
			flows = append(flows, newflows...)
			continue
		}

		inps := inPortsFromOutPorts(outports, ctx)
		outports = outports[0:0]
		for _, inp := range inps {
			newflows, newoutports, err := genFlowsForPathlet(p, inp,
				add.Path.Priority, ctx)
			if err != nil {
				return err
			}
			outports = append(outports, newoutports...)
			flows = append(flows, newflows...)
		}
	}

	uniqueFlows := make([]nom.FlowEntry, 0, len(flows))
nextFlow:
	for i, fi := range flows {
		for j, fj := range flows {
			if i == j {
				continue
			}

			if fj.Subsumes(fi) {
				continue nextFlow
			}

			// TODO(soheil): check for subsumption and merge flows if possible.
			if j < i && fj.Equals(fi) {
				continue nextFlow
			}
		}
		uniqueFlows = append(uniqueFlows, fi)
	}
	addFlowEntriesForPath(add.Subscriber, add.Path, uniqueFlows, ctx)
	return nil
}

func (h addHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return centralizedMap
}

func genFlowsForPathlet(p nom.Pathlet, inport nom.UID, priority uint16,
	ctx bh.RcvContext) (flows []nom.FlowEntry, outports []nom.UID, err error) {

	fwdnps := forwardNodes(p.Actions)
	for _, ports := range fwdnps {
		outports = append(outports, ports...)
	}
	floodns := floodNodes(p.Actions)
	for n, p := range floodns {
		outports = append(outports, outPortsFromFloodNode(n, p, ctx)...)
	}

	port, matchHasPort := p.Match.InPort()
	if matchHasPort {
		if inport != nom.Nil && inport != nom.UID(port) {
			return nil, nil, fmt.Errorf("path: two different inports %v and %v",
				inport, port)
		}
		inport = nom.UID(port)
	}

	m := p.Match
	if inport != nom.Nil && !matchHasPort {
		m = p.Match.Clone()
		m.Fields = append(m.Fields, nom.InPort(inport))
	}

	noinMatch := p.Match.Clone()
	for f := range m.Fields {
		if _, ok := m.Fields[f].(nom.InPort); ok {
			noinMatch.Fields = append(noinMatch.Fields[:f], m.Fields[f+1:]...)
			break
		}
	}

	nofwdActions := make([]nom.Action, 0, len(p.Actions))
	for _, a := range p.Actions {
		switch a.(type) {
		case nom.ActionForward, nom.ActionFlood:
			continue
		default:
			nofwdActions = append(nofwdActions, a)
		}
	}

	var innodes []nom.UID
	if inport != nom.Nil {
		innodes = []nom.UID{nom.NodeFromPortUID(inport)}
	} else {
		innodes = discovery.NodesCentralized(ctx)
	}

	for _, inn := range innodes {
		for _, outp := range outports {
			outn := nom.NodeFromPortUID(outp)
			sps, l := discovery.ShortestPathCentralized(inn, outn, ctx)
			if l < 0 {
				// TODO(soheil): maybe just log this and continue installing other
				// flows.
				return nil, nil, fmt.Errorf("path: no path found from %v to %v", inport,
					outp)
			}

			if l == 0 {
				m := noinMatch.Clone()
				if inport != nom.Nil {
					m.Fields = append(m.Fields, nom.InPort(inport))
				}
				flow := nom.FlowEntry{
					Node:     outn,
					Match:    m,
					Actions:  p.Actions,
					Priority: priority,
				}
				flows = append(flows, flow)
				continue
			}

			// TODO(soheil): maybe install multiple paths.
			lastInPort := inport
			for _, link := range sps[0] {
				m := noinMatch.Clone()
				if lastInPort != nom.Nil {
					m.Fields = append(m.Fields, nom.InPort(lastInPort))
				}

				var a []nom.Action
				a = append(a, nofwdActions...)
				a = append(a, nom.ActionForward{Ports: []nom.UID{link.From}})

				flow := nom.FlowEntry{
					Node:     nom.NodeFromPortUID(link.From),
					Match:    m,
					Actions:  a,
					Priority: priority,
				}
				flows = append(flows, flow)

				lastInPort = link.To
			}

			m := noinMatch.Clone()
			if lastInPort != nom.Nil {
				m.Fields = append(m.Fields, nom.InPort(lastInPort))
			}
			flow := nom.FlowEntry{
				Node:     outn,
				Match:    m,
				Actions:  p.Actions,
				Priority: priority,
			}
			flows = append(flows, flow)
		}
	}
	return flows, outports, nil
}

func forwardNodes(actions []nom.Action) (nodeToPorts map[nom.UID][]nom.UID) {
	nodeToPorts = make(map[nom.UID][]nom.UID)
	for _, a := range actions {
		switch f := a.(type) {
		case nom.ActionForward:
			for _, p := range f.Ports {
				fnid, _ := nom.ParsePortUID(p)
				ports := nodeToPorts[fnid.UID()]
				ports = append(ports, p)
				nodeToPorts[fnid.UID()] = ports
			}
		}
	}
	return nodeToPorts
}

func floodNodes(actions []nom.Action) (nodes map[nom.UID]nom.UID) {
	fn := make(map[nom.UID]nom.UID)
	for _, a := range actions {
		switch f := a.(type) {
		case nom.ActionFlood:
			nid := nom.NodeFromPortUID(f.InPort)
			fn[nid] = f.InPort
		}
	}
	return fn
}
