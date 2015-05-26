package kandoo

import (
	bh "github.com/kandoo/beehive"
)

// Implements the map function for local handlers.
type Local struct{}

func (h Local) Map(msg bh.Msg, ctx bh.MapContext) bh.MappedCells {
	return ctx.LocalMappedCells()
}
