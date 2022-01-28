package server

import (
	"io"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"x.x/x/deweb/crypt"
	"x.x/x/deweb/lib"
)

func init() {
	print(uuid.New().String())
}

func handleconn(conn net.Conn) {
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	b, _ := io.ReadAll(conn)
	var response lib.TransportStruct
	err := crypt.Decrypt(b, &response)
	if err != nil {
		panic(err)
	}
	method_split := strings.Split(response.Method, "/")
	// v1  0   ping
	// [0] [1] [2]
	if len(method_split) != 3 {
		panic("len(response.Method) != 3")
	}
	switch response.Method {
	case "v1/0/ping":

	default:
		print(response.Method + " is not supported\n")
	}

}
