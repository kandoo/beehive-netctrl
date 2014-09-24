package nom

// PacketIn messages are emitted when a packet is forwarded to the controller.
type PacketIn struct {
	Node   UID
	InPort UID
	Packet Packet
}

// PacketOut messages are emitted to send a packet out of a port.
type PacketOut struct {
	Node    UID
	OutPort UID
	InPort  UID
	Packet  Packet
}

// Packet is simply the packet data.
type Packet []byte
