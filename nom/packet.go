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

// PacketBufferID represents a packet buffered in the switch.
type PacketBufferID uint32
