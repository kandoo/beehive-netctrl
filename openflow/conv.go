package openflow

import (
	"strconv"

	"github.com/soheilhy/beehive-netctrl/nom"
)

func datapathIDToNodeID(datapathID uint64) nom.NodeID {
	return nom.NodeID(strconv.FormatUint(datapathID, 16))
}

func nodeIDToDatapathID(nodeID nom.NodeID) (uint64, error) {
	return strconv.ParseUint(string(nodeID), 16, 64)
}

func portNoToPortID(portNo uint32) nom.PortID {
	return nom.PortID(strconv.FormatUint(uint64(portNo), 10))
}

func portIDToPortNo(portID nom.PortID) (uint32, error) {
	no, err := strconv.ParseUint(string(portID), 10, 32)
	return uint32(no), err
}
