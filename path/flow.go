package path

import (
	"fmt"
	"strconv"
	"time"

	bh "github.com/kandoo/beehive"
	"github.com/kandoo/beehive-netctrl/nom"
	"github.com/kandoo/beehive/Godeps/_workspace/src/github.com/golang/glog"
)

type pathAndFlows struct {
	Subscriber bh.AppCellKey
	Path       nom.Path
	Flows      []flowAndStatus
	Installed  int
	Timestamp  time.Time
}

type flowAndStatus struct {
	Flow      nom.FlowEntry
	Installed bool
}

func addFlowEntriesForPath(sub bh.AppCellKey, path nom.Path,
	flows []nom.FlowEntry, ctx bh.RcvContext) {

	fs := make([]flowAndStatus, 0, len(flows))
	path.ID = strconv.FormatUint(reservePathID(ctx), 16)
	for i := range flows {
		flows[i].ID = path.ID
		fs = append(fs, flowAndStatus{Flow: flows[i]})
	}

	pf := pathAndFlows{
		Subscriber: sub,
		Path:       path,
		Flows:      fs,
		Timestamp:  time.Now(),
	}
	d := ctx.Dict(dictPath)
	if err := d.Put(path.ID, pf); err != nil {
		glog.Fatalf("error in storing path entry: %v", err)
	}

	ack := centralizedAppCellKey(ctx.App())
	for _, f := range flows {
		addf := nom.AddFlowEntry{
			Flow:       f,
			Subscriber: ack,
		}
		ctx.Emit(addf)
	}
}

func confirmFlowEntryForPath(flow nom.FlowEntry, ctx bh.RcvContext) error {
	d := ctx.Dict(dictPath)

	v, err := d.Get(flow.ID)
	if err != nil {
		return fmt.Errorf("path: flow not found: %v", err)
	}

	pf := v.(pathAndFlows)

	for i := range pf.Flows {
		if pf.Flows[i].Flow.Equals(flow) {
			if pf.Flows[i].Installed {
				return fmt.Errorf("%v is already installed", flow)
			}
			pf.Flows[i].Installed = true
			pf.Installed++
			break
		}
	}

	if pf.Installed == len(pf.Flows) {
		ctx.SendToCell(nom.PathAdded{Path: pf.Path}, pf.Subscriber.App,
			pf.Subscriber.Cell())
	}
	return d.Put(flow.ID, pf)
}

func delFlowEntryFromPath(flow nom.FlowEntry, ctx bh.RcvContext) error {
	// TODO(soheil): Add deleted flow entries.
	panic("TODO not implemented yet")
}

type flowHandler struct{}

func (h flowHandler) Rcv(msg bh.Msg, ctx bh.RcvContext) error {
	switch data := msg.Data().(type) {
	case nom.FlowEntryAdded:
		return confirmFlowEntryForPath(nom.FlowEntry(data.Flow), ctx)
	case nom.FlowEntryDeleted:
		return delFlowEntryFromPath(nom.FlowEntry(data.Flow), ctx)
	}
	return fmt.Errorf("flowHandler: unsupported message %v", msg.Type())
}

func (h flowHandler) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return nil
}
