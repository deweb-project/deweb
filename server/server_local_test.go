package server_test

import (
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"x.x/x/deweb/crypt"
	"x.x/x/deweb/justdb"
	"x.x/x/deweb/lib"
	"x.x/x/deweb/server"
	"x.x/x/deweb/transport"
)

func TestConnection(t *testing.T) {
	justdb.Setup("/dev/shm/testconn")
	crypt.LoadSelfKey()
	alice_server := server.NewServerLocal(50001)
	bob_server := server.NewServerLocal(50002)
	go func() {
		log.Println("[alice] Starting...")
		log.Println("[alice]", alice_server.Start())
		t.Fail()
	}()
	go func() {
		log.Println("[bob] Starting...")
		log.Println("[bob]", bob_server.Start())
		t.Fail()
	}()
	time.Sleep(time.Second / 2)
	aliceconn := transport.GetTransportLocal(alice_server.DEID)
	id := uuid.New().String()
	response := lib.TransportStruct{
		ID:          []byte(id),
		Nonce:       id,
		Source:      lib.GetSelfID().ID,
		Destination: alice_server.DEID.String(),
		Method:      "v1/0/ping",
		Data:        "ping",
	}
	response.OUTAttachPublicKey()
	response.OUTAttachSignature()
	//response.INVerifyMessage()
	err := aliceconn.Send(alice_server.DEID, aliceconn.Key, response)
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	time.Sleep(time.Second * 5) // wait for async jobs to not panic
	//t.Fail()
}
