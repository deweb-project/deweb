package server

import (
	"log"
	"net"
)

func NewServerLocal() Server {
	return Server{
		Start: func() error {
			log.Println("Listening on port :51337")
			l, err := net.Listen("tcp", ":51337")
			if err != nil {
				log.Fatal(err)
			}
			for {
				conn, err := l.Accept()
				if err != nil {
					log.Println(err)
					continue
				}
				go handleconn(conn)
			}
		},
	}
}
