package nom

import "encoding/gob"

// DriverRole is the role of a driver for a node. The values can be either
// default, slave or master. Only one driver can be the master, but we can
// have multiple slaves and multiple defaults.
type DriverRole uint8

// Valid values for DriverRole.
const (
	DriverRoleDefault DriverRole = iota
	DriverRoleSlave              = iota
	DriverRoleMaster             = iota
)

func init() {
	gob.Register(DriverRole(0))
}
