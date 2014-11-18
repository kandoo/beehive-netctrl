package path

import (
	bh "github.com/kandoo/beehive"
)

func reserveFlowID(ctx bh.RcvContext, cnt int) uint64 {
	d := ctx.Dict(dictID)
	var id uint64
	d.GetGob("flow", &id)
	id += uint64(cnt)
	d.PutGob("flow", &id)
	return id - uint64(cnt)
}

func reservePathID(ctx bh.RcvContext) uint64 {
	d := ctx.Dict(dictID)
	var id uint64
	d.GetGob("path", &id)
	id++
	d.PutGob("path", &id)
	return id - 1
}
