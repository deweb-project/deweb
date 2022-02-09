package server

import (
	"x.x/x/deweb/lib"
	"x.x/x/deweb/transport"
)

type Server struct {
	Start      func() error
	DEID       lib.DEID
	Connection transport.Connection
}
