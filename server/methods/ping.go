package methods

import (
	"github.com/google/uuid"
	"x.x/x/deweb/lib"
)

func HandlePingV1(x lib.TransportStruct) {
	// pong is response, in this case we don't send anything back.
	//log.Println(x.Method)
	if x.Data == "pong" {
		return
	}
	id := uuid.New().String()
	response := lib.TransportStruct{
		ID:          []byte(id),
		Nonce:       id,
		Source:      lib.GetSelfID().ID,
		Destination: x.Source,
		Method:      x.Method,
		Data:        "pong",
	}
	response.OUTAttachPublicKey()
	response.OUTAttachSignature()
	//response.INVerifyMessage()
	lib.QueueTask(x)
}
