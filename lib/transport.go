package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/google/uuid"
	"x.x/x/deweb/crypt"
)

// This TransportStruct
type TransportStruct struct {
	ID          []byte // Nonce but []byte, for JustDB
	Nonce       string // Random String ~
	Source      string // GetSelfID().ID
	Destination string // GetSelfID().ID
	Method      string // For example v1/0/ping
	// Where v1 is major version - bump when breaking change occured
	// 0 is a minor version.
	// ping is a command, commands should not use `/` as separator - use `-` instead
	// v1/0/chat-messages-send
	Data      string `json:",omitempty"`
	DataBytes []byte `json:",omitempty"`
	PublicKey string `json:",omitempty"` // outcoming pubkey.
	Signature string // GenerateSignatureString(.this)
	Tries     int    // How many times the event tried to be delivered.
}

func (ts *TransportStruct) OUTInitNonce() {
	nonce := uuid.New().String()
	ts.ID = []byte(nonce)
	ts.Nonce = nonce
}

func (ts *TransportStruct) String() string {
	b, _ := json.Marshal(ts)
	return string(b)
}

func (ts *TransportStruct) OUTAttachPublicKey() {
	if ts.PublicKey != "" {
		return
	}
	var err error
	ts.PublicKey, err = crypt.GetKey().GetArmoredPublicKey()
	if err != nil {
		panic(err)
	}
}

func (ts *TransportStruct) OUTAttachSignature() {
	keyring, err := crypto.NewKeyRing(crypt.GetKey())
	if err != nil {
		panic(err)
	}
	msg := crypto.NewPlainMessageFromString(ts.GenerateSignatureString())
	pgpsign, err := keyring.SignDetached(msg)
	if err != nil {
		panic(err)
	}
	ts.Signature, err = pgpsign.GetArmored()
	if err != nil {
		panic(err)
	}
}

func (ts *TransportStruct) GenerateSignatureString() string {
	return "[" + ts.Source + "]&[" + ts.Destination + "]&[" + ts.Nonce + "]&[" + ts.HashContent() + "]"
}

func (ts *TransportStruct) HashContent() string {
	//type TransportStruct struct {
	//	ID          []byte // Nonce but []byte, for JustDB
	//	Source      string // GetSelfID().ID
	//	Destination string // GetSelfID().ID
	//	Nonce       string // Random String ~
	//	Method      string // For example v1/0/ping
	//	// Where v1 is major version - bump when breaking change occured
	//	// 0 is a minor version.
	//	// ping is a command, commands should not use `/` as separator - use `-` instead
	//	// v1/0/chat-messages-send
	//	Data      string `json:",omitempty"`
	//	DataBytes []byte `json:",omitempty"`
	//	PublicKey string `json:",omitempty"`
	//	Signature string // GenerateSignatureString(.this)
	//}
	tohash := fmt.Sprintf(`ID: %x,
Source: %s,
Destination: %s,
Nonce: %s,
Method: %s,
Data: %s,
DataBytes: %x`, ts.ID, ts.Source, ts.Destination, ts.Nonce, ts.Method, ts.Data, ts.DataBytes)
	return crypt.SHA512(tohash)
}

// true - message ok
func (ts *TransportStruct) INVerifyMessage() bool {
	//TODO: disable this, full trust.
	//return true
	proto := strings.Split(ts.Source, ":")[0]
	if proto == "local" || proto == "dummyproto" {
		return true
	}
	//log.Println(ts.Source)
	msg := crypto.NewPlainMessageFromString(ts.GenerateSignatureString())
	signature, err := crypto.NewPGPSignatureFromArmored(ts.Signature)
	if err != nil {
		log.Println(err, ts.Signature)
		return false
	}
	pubkey, err := crypto.NewKeyFromArmored(ts.PublicKey)
	if err != nil {
		log.Println(err)
		return false
	}

	signkey, err := crypto.NewKeyRing(pubkey)
	if err != nil {
		log.Println(err)
		return false
	}

	err = signkey.VerifyDetached(msg, signature, crypto.GetUnixTime())
	if err != nil {
		log.Println(err)
		return false
	}

	deid, err := ParseDEID(ts.Source)
	if err != nil {
		log.Println(err)
		return false
	}
	if deid.Key != pubkey.GetFingerprint() {
		log.Println(deid.Key, "(deid.Key) != ", pubkey.GetFingerprint()+"(pubkey.GetFingerprint())")
		return false
	}
	return true
}
