package controller

import (
	"encoding/json"

	"github.com/soheilhy/beehive-netctrl/nom"
)

const (
	nodeDriversDict = "N"
)

type nodeDrivers struct {
	Node    nom.Node
	Drivers []nom.Driver
	Ports   nom.Ports
}

func (nd *nodeDrivers) UID() nom.UID {
	return nd.Node.UID()
}

func (nd *nodeDrivers) GoDecode(b []byte) error {
	return nom.ObjGoDecode(nd, b)
}

func (nd *nodeDrivers) GoEncode() ([]byte, error) {
	return nom.ObjGoEncode(nd)
}

func (nd *nodeDrivers) JSONDecode(b []byte) error {
	return json.Unmarshal(b, nd)
}

func (nd *nodeDrivers) JSONEncode() ([]byte, error) {
	return json.Marshal(nd)
}

func (nd nodeDrivers) hasDriver(d nom.Driver) bool {
	for _, e := range nd.Drivers {
		if e == d {
			return true
		}
	}

	return false
}

func (nd *nodeDrivers) removeDriver(d nom.Driver) bool {
	for i, e := range nd.Drivers {
		if e == d {
			nd.Drivers = append(nd.Drivers[:i], nd.Drivers[i+1:]...)
			return true
		}
	}
	return false
}

func (nd *nodeDrivers) master() nom.Driver {
	// FIXME(soheil)
	return nd.Drivers[0]
}
