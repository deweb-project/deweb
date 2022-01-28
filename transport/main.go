package transport

import "errors"

var ConnectionIsUnreachable = errors.New("connection: host is unreachable")

type Connection struct {
	Protocol    string
	Destination string

	EstabilishConnection func() error              // This function is being called before first connection attampt, and everytime the connection is lost.
	Send                 func(x interface{}) error // Send something to the target, error = not received.
	Ping                 func() error
}
