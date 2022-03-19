package methods

import (
	"log"

	"x.x/x/deweb/justdb"
	"x.x/x/deweb/lib"
)

func HandleChatInviteV1(x lib.TransportStruct) {
	var conv = lib.Conversation{
		ID:       x.DataBytes,
		Host:     x.Source,
		DEID:     []string{x.Source},
		RoomName: "Unknown Room",
	}
	if !lib.IsAdminInConversation(x.Source, lib.Conversation{ID: x.DataBytes}) {
		log.Println("got invited by non-admin, ignoring beacuse we won't get accepted.")
		return
	}
	justdb.Write(&conv)
}
