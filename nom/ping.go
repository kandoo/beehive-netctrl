package nom

import "encoding/gob"

// Ping represents a ping message sent to a driver. It is used to
// health-check the driver.
type Ping struct{}

// Pong represents a reply to a ping message sent by the driver.
type Pong struct{}

func init() {
	gob.Register(Ping{})
	gob.Register(Pong{})
}
