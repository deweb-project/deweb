package lib

import (
	"fmt"
	"log"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
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
	PublicKey string `json:",omitempty"`
	Signature string // GenerateSignatureString(.this)
	Tries     int    // How many times the event tried to be delivered.
}

func (ts *TransportStruct) OUTAttachPublicKey() {
	if ts.PublicKey != "" {
		return
	}
	var err error
	ts.PublicKey, err = crypt.Key.GetArmoredPublicKey()
	if err != nil {
		panic(err)
	}
}

func (ts *TransportStruct) OUTAttachSignature() {
	keyring, err := crypto.NewKeyRing(crypt.Key)
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
	msg := crypto.NewPlainMessageFromString(ts.GenerateSignatureString())
	signature, err := crypto.NewPGPSignatureFromArmored(ts.Signature)
	if err != nil {
		log.Println(err)
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
