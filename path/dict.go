package path

import (
	bh "github.com/kandoo/beehive"
)

const (
	centralizedD = "D"
	centralizedK = "0"

	dictFlow = "FlowDict"
	dictPath = "PathDict"
	dictID   = "IDDict"
)

var centralizedMap = bh.MappedCells{{Dict: centralizedD, Key: centralizedK}}

func centralizedAppCellKey(app string) bh.AppCellKey {
	return bh.AppCellKey{
		App:  app,
		Dict: centralizedD,
		Key:  centralizedK,
	}
}
