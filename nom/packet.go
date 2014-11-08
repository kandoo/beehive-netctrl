package nom

// Special ports.
const (
	PortFlood UID = "Ports.PortBcast"
	PortAll   UID = "Ports.PortAll"
)

// PacketIn messages are emitted when a packet is forwarded to the controller.
type PacketIn struct {
	Node     UID
	InPort   UID
	BufferID PacketBufferID
	Packet   Packet
}

// PacketOut messages are emitted to send a packet out of a port.
type PacketOut struct {
	Node     UID
	InPort   UID
	BufferID PacketBufferID
	Packet   Packet
	Actions  []Action
}

// Packet is simply the packet data.
type Packet []byte

// DstMAC returns the destination MAC address from the ethernet header.
func (p Packet) DstMAC() MACAddr {
	return MACAddr{p[0], p[1], p[2], p[3], p[4], p[5]}
}

// SrcMAC returns the source MAC address from the ethernet header.
func (p Packet) SrcMAC() MACAddr {
	return MACAddr{p[6], p[7], p[8], p[9], p[10], p[11]}
}

// TODO(soheil): add code to parse ip addresses and tcp ports.

// PacketBufferID represents a packet buffered in the switch.
type PacketBufferID uint32
