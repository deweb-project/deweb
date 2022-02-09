package methods

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"time"

	"x.x/x/deweb/justdb"
	"x.x/x/deweb/lib"
)

type Conversation struct {
	ID           []byte // Random ID
	Deid         []string
	RoomName     string
	MessageIndex int
}
type Message struct {
	ID          []byte // Conversation.ID + Incremental ID
	From        string
	ContentType string
	Content     []byte
	Sent        time.Time
	Received    time.Time
}

type ConversationList struct {
	ID             []byte
	ConversationID [][]byte
}

var conversationlist = ConversationList{
	ID: []byte("account1"),
}

func HandleMessageV1(x lib.TransportStruct) {
	spl := strings.Split(x.Data, ":") // <Content-Type>:<ConversationID>:time.Now().UnixNano()Message.Sent
	if len(spl) != 3 {
		log.Println(x.Data, " len != 3")
		return
	}
	sent, err := strconv.ParseInt(spl[2], 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	var conv = Conversation{
		ID: []byte(spl[1]),
	}

	justdb.Read(&conv)
	conv.MessageIndex++
	justdb.Write(&conv)
	msg := Message{
		ID:          append(conv.ID, []byte("|"+strconv.Itoa(conv.MessageIndex))...),
		From:        x.Source,
		ContentType: spl[0],
		Content:     x.DataBytes,
		Sent:        time.Unix(0, sent),
		Received:    time.Now(),
	}
	justdb.Write(msg)
}

func isInConversationIndex(id []byte) bool {
	for _, v := range conversationlist.ConversationID {
		if bytes.Compare(v, id) == 0 {
			return true
		}
	}
	return false
}
