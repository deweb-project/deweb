package transport

import (
	"errors"

	"x.x/x/deweb/lib"
)

var ConnectionIsUnreachable = errors.New("connection: host is unreachable")

type Connection struct {
	Protocol             string
	Destination          string
	Key                  string
	EstabilishConnection func() error                                                 // This function is being called before first connection attampt, and everytime the connection is lost.
	Send                 func(deid lib.DEID, key string, x lib.TransportStruct) error // Send something to the target, error = not received.
	Ping                 func() error
	GetKey               func(deid lib.DEID) (key string, err error)
}
