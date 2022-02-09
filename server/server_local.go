package server

import (
	"log"
	"net"
	"strconv"

	"x.x/x/deweb/crypt"
	"x.x/x/deweb/lib"
	"x.x/x/deweb/transport"
)

func NewServerLocal(port int) Server {
	deidstring := "local:127.0.0.1:" + strconv.Itoa(port) + "[key=" + crypt.GetKey().GetFingerprint() + "]"
	deid, err := lib.ParseDEID(deidstring)
	if err != nil {
		panic(err)
	}
	return Server{
		DEID: deid,
		Start: func() error {
			log.Println("Listening on port :" + strconv.Itoa(port))
			l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
			if err != nil {
				log.Fatal(err)
			}
			for {
				conn, err := l.Accept()
				if err != nil {
					log.Println(err)
					continue
				}
				go transport.Handleconn(conn, false)
			}
		},
	}
}
