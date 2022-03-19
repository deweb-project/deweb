package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"x.x/x/deweb/justdb"
)

type Conversation struct {
	ID           []byte   // Random ID (32 chars) + | + self.deid
	Host         string   // self.deid
	DEID         []string // members
	RoomName     string   // Whatever
	MessageIndex int      // Internal use
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

// target - Conversation.ID
func SendMessage(target []byte, text string, data []byte) {
	var packet TransportStruct
	// List of content-type:
	// - text/v0
	// - notice/v0
	packet.Data = "text/v0:" + string(target) + ":" + fmt.Sprint(time.Now().UnixNano()) + ":" + text // <Content-Type>:<ConversationID>:time.Now().UnixNano()Message.Sent:fallback text
	packet.Source = GetSelfID().ID
	packet.DataBytes = data
	packet.Method = "v1/0/message"
	// The packet is ready, let's replicate it for each DEID in Conversation

	var conversation Conversation
	conversation.ID = target
	justdb.Read(&conversation)
	for i := range conversation.DEID {
		forkedpacket := packet
		forkedpacket.Destination = conversation.DEID[i]
		//forkedpacket.OUTInitNonce()
		QueueTask(forkedpacket)
	}
}

// Internal use only - using ParseMessage will result in a write call to justdb, every time it is called
// Avoid using externally.
func HandleNewMessage(x TransportStruct) Message {
	spl := strings.Split(x.Data, ":") // <Content-Type>:<ConversationID>:time.Now().UnixNano()Message.Sent:fallback text
	if len(spl) != 3 {
		log.Println(x.Data, " len != 3")
		return Message{}
	}
	sent, err := strconv.ParseInt(spl[2], 10, 64)
	if err != nil {
		log.Println(err)
		return Message{}
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
	return msg
}

func LoadConversationList() {
	justdb.Read(&conversationlist)
}

var conversationlist = ConversationList{
	ID: []byte("account1"),
}

func isInConversationIndex(id []byte) bool {
	for _, v := range conversationlist.ConversationID {
		if bytes.Compare(v, id) == 0 {
			return true
		}
	}
	return false
}

// Now let's think for a bit how can we handle conversation
// First of all we need a method to actually create a conversation
// This should leave us with a empty conversation, that we can configure.
// Then we should be able to add other users to the to the conversation,
//and sign the new list of users with our (creator's) key, and broadcast
//it to every existing user. That way we don't have to sync old events
//with devices that were left offline for long time, and did not receive
//all events that happened.
// Speaking of which - `Conversation` struct will be synced between all
//participants in a conversation in a one to many way - the creator have
//the final word.
// However, not now because this task will be too much time consuming,
//it will be possible to 'fork' a group. If an admin is offline, ignore
//user's requests, is pro-censorship, then users should be able to easily
//fork the group, where one user sends a fork request, and others are left
//with an option to accept it or ignore it.

// CreateConversation
// name - how do you want to call this conversation?
func CreateConversation(name string) Conversation {
	var id = randomString(32) + "|" + GetSelfID().ID
	var conv = Conversation{
		ID:       []byte(id),
		RoomName: id,
	}
	err := justdb.Write(&conv)
	if err != nil {
		panic(err)
	}
	return conv
}

// Send an invitation request to some guy.
func SendInvitation(user DEID, conv Conversation) {
	//var inviteString = "<Conversation.ID>"
	var request = TransportStruct{
		Method:      "v1/0/chat-invite",
		DataBytes:   conv.ID,
		Destination: user.String(),
	}
	QueueTask(request)
}
func IsAdminInConversation(user string, conv Conversation) bool {
	cs := string(conv.ID)
	return strings.HasSuffix(cs, user)
}

type Invitation struct {
	ID     []byte // []byte(DEID)
	DEID   string
	ConvID []byte
}

func AcceptInvitation(inv Invitation) {
	var request = TransportStruct{
		Method:      "v1/0/chat-invite-accept",
		DataBytes:   inv.ID,
		Destination: inv.DEID,
	}
	QueueTask(request)
}

func IgnoreInvitation(inv Invitation) {
	justdb.Delete(inv)
}
func AnnounceConversation(conv Conversation) {
	// type Conversation struct {
	// 	ID           []byte // Random ID (32 chars) + | + self.deid
	// 	Host         string
	// 	DEID         []string
	// 	RoomName     string
	// 	MessageIndex int
	// }

	c, err := json.Marshal(conv)
	if err != nil {
		log.Fatal(err)
	}
	for i := range conv.DEID {
		var req = TransportStruct{
			Method:      "v1/0/chat-conversation-update",
			DataBytes:   c,
			Destination: conv.DEID[i],
		}
		QueueTask(req)
	}
}
func BanUser(conv Conversation, targetDEID string) {
	for i := range conv.DEID {
		if targetDEID == conv.DEID[i] {
			conv.DEID = append(conv.DEID[:i], conv.DEID[i+1:]...)
		}
	}
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
