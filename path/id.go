package path

import (
	bh "github.com/kandoo/beehive"
)

func reserveFlowID(ctx bh.RcvContext, cnt int) uint64 {
	d := ctx.Dict(dictID)
	var id uint64
	if v, err := d.Get("flow"); err == nil {
		id = v.(uint64)
	}
	id += uint64(cnt)
	d.Put("flow", id)
	return id - uint64(cnt)
}

func reservePathID(ctx bh.RcvContext) uint64 {
	d := ctx.Dict(dictID)
	var id uint64
	if v, err := d.Get("path"); err == nil {
		id = v.(uint64)
	}
	id++
	d.Put("path", id)
	return id - 1
}
