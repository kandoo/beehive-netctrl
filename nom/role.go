package nom

import "encoding/gob"

// DriverRole is the role of a driver for a node. The values can be either
// default, slave or master. Only one driver can be the master, but we can
// have multiple slaves and multiple defaults.
type DriverRole uint8

// Valid values for DriverRole.
const (
	DriverRoleDefault DriverRole = iota
	DriverRoleSlave
	DriverRoleMaster
)

// ChangeDriverRole is emitted to instruct a driver to change its role for a
// node.
type ChangeDriverRole struct {
	Node       UID        // Tho node ID.
	Role       DriverRole // The requested new role.
	Generation uint64     // The generation of role request.
}

// DriverRoleUpdate is a message emitted when a driver's role is changed or
// an update (no necessarily a change) is recevied for a node.
type DriverRoleUpdate struct {
	Node       UID
	Driver     Driver
	Generation uint64
}

func init() {
	gob.Register(ChangeDriverRole{})
	gob.Register(DriverRole(0))
	gob.Register(DriverRoleUpdate{})
}
