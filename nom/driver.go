package nom

import (
	"encoding/json"
	"fmt"
)

// Driver represents the Bee that communicates with the actual networking
// element and adapts the southbound protocol for using in NOM.
type Driver struct {
	BeeID uint64
	Role  DriverRole
}

func (d Driver) String() string {
	return fmt.Sprintf("%016X (role=%v)", d.BeeID, d.Role)
}

// GoDecode decodes the driver from the byte slice using Gob.
func (d *Driver) GoDecode(b []byte) error {
	return ObjGoDecode(d, b)
}

// GoEncode encodes the driver into a byte slice using Gob.
func (d *Driver) GoEncode() ([]byte, error) {
	return ObjGoEncode(d)
}

// JSONDecode decodes the driver from a byte slice using JSON.
func (d *Driver) JSONDecode(b []byte) error {
	return json.Unmarshal(b, d)
}

// JSONEncode encodes the driver into a byte slice using JSON.
func (d *Driver) JSONEncode() ([]byte, error) {
	return json.Marshal(d)
}

type Drivers []Driver

// GoDecode decodes the drivers from the byte slice using Gob.
func (d *Drivers) GoDecode(b []byte) error {
	return ObjGoDecode(d, b)
}

// GoEncode encodes the drivers into a byte slice using Gob.
func (d *Drivers) GoEncode() ([]byte, error) {
	return ObjGoEncode(d)
}

// JSONDecode decodes the drivers from a byte slice using JSON.
func (d *Drivers) JSONDecode(b []byte) error {
	return json.Unmarshal(b, d)
}

// JSONEncode encodes the drivers into a byte slice using JSON.
func (d *Drivers) JSONEncode() ([]byte, error) {
	return json.Marshal(d)
}
