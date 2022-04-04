package transport

import (
	"io"
	"log"
	"net"
	"strings"
	"time"

	"x.x/x/deweb/crypt"
	"x.x/x/deweb/lib"
	"x.x/x/deweb/server/methods"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func Handleconn(conn net.Conn, readreply bool) {
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	// read the high level endpoint
	data := ""
	b, _ := io.ReadAll(conn)
	//if err != "panic: read tcp 127.0.0.1:50001->127.0.0.1:52540: i/o timeout" {
	//	panic(err)
	//}
	data = string(b)
	if strings.HasPrefix(data, "plaintext/init/v1\n") {
		log.Println("Plaintext prefix!")
		pk, err := crypt.GetKey().GetArmoredPublicKey()
		if err != nil {
			panic(err)
		}
		conn.Write([]byte(lib.GetSelfID().ID + "\n" + pk))
		conn.Close()
		return
	}
	log.Println("data: len(b):", len(b))
	if len(b) == 0 {
		return
	}
	//if len("plaintext/init/v1\n") > len(data) {
	//	log.Println("breaking")
	//}
	//b_part2, err := ioutil.ReadAll(conn)
	//if err != nil {
	//	log.Println(err)
	//	conn.Close()
	//	return
	//}
	//b := append([]byte(data), b_part2...)
	var response lib.TransportStruct
	log.Println("crypt.Decrypt(b, &response)")
	err := crypt.Decrypt(b, &response)
	log.Println("crypt.Decrypt(b, &response): end")
	if !response.INVerifyMessage() {
		log.Println("INVerifyMessage() failed, bad actor possible")
		conn.Close()
		return
	}
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Println("RESPONSE:", response)
	method_split := strings.Split(response.Method, "/")
	// v1  0   ping
	// [0] [1] [2]
	if len(method_split) != 3 {
		panic("len(response.Method) != 3")
	}

	switch response.Method {
	case "v1/0/ping":
		methods.HandlePingV1(response)
	case "v1/0/message":
		methods.HandleMessageV1(response)
	case "v1/0/chat-invite":
		methods.HandleChatInviteV1(response)
	case "v1/0/chat-invite-accept":
		log.Println(response, "TODO")
	default:
		log.Println(response.Method + " is not supported\n")
	}
	if readreply {
		task := lib.GetTask(response.Source)
		data := crypt.Encrypt(&task, response.PublicKey)
		conn.Write(data)
	}
	conn.Close()
}
