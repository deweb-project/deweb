package lib_test

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/google/uuid"
	"x.x/x/deweb/crypt"
	"x.x/x/deweb/lib"
)

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func generateKey() *crypto.Key {
	rsaKey, err := helper.GenerateKey("na", "na@na.na", []byte{}, "rsa", 1024) // 1024 just for testing
	if err != nil {
		panic(err)
	}
	key, err := crypto.NewKeyFromArmored(rsaKey)
	if err != nil {
		panic(err)
	}
	return key
}

func TestTransportEncryptionCorrect(t *testing.T) {
	// Correct
	for i := 0; i < 75; i++ {

		nonce := uuid.New().String()
		key_receiver := generateKey()
		crypt.Key = key_receiver
		//crypt.KeyUnlockDefault()
		deid_receiver := lib.GetSelfID().ID

		key_sender := generateKey()
		crypt.Key = key_sender
		//crypt.KeyUnlockDefault()
		deid_sender := lib.GetSelfID().ID

		var message = lib.TransportStruct{
			ID:          []byte(nonce),
			Source:      deid_sender,
			Destination: deid_receiver,
			Nonce:       nonce,
			Method:      "v1/0/ping",
			Data:        randomString(2048),
			DataBytes:   []byte(randomString(2048)),
		}
		message.OUTAttachPublicKey()
		message.OUTAttachSignature()
		// message signed. Now let's verify it.
		key_receiver = generateKey()
		crypt.Key = key_receiver
		//crypt.KeyUnlockDefault()
		ok := message.INVerifyMessage()
		if ok != true {
			t.Error("message.INVerifyMessage(): failed")
			t.Fail()
		}
	}
}

func TestTransportEncryptionIncorrect(t *testing.T) {
	// Incorrect - things changed
	for i := 0; i < 1; i++ {
		nonce := uuid.New().String()
		key_receiver := generateKey()
		crypt.Key = key_receiver
		//crypt.KeyUnlockDefault()
		deid_receiver := lib.GetSelfID().ID

		key_sender := generateKey()
		crypt.Key = key_sender
		//crypt.KeyUnlockDefault()
		deid_sender := lib.GetSelfID().ID

		var message = lib.TransportStruct{
			ID:          []byte(nonce),
			Source:      deid_sender,
			Destination: deid_receiver,
			Nonce:       nonce,
			Method:      "v1/0/ping",
			Data:        randomString(2048),
			DataBytes:   []byte(randomString(2048)),
		}

		message.OUTAttachPublicKey()
		message.OUTAttachSignature()
		// message signed. Now let's verify it.
		switch i % 5 {
		case 0:
			message.Data = randomString(2048)
		case 1:
			message.DataBytes = []byte(randomString(2048))
		case 2:
			message.Data = randomString(2048)
			message.DataBytes = []byte(message.Data)
		case 3:
			message.Source = deid_receiver
			message.Destination = deid_sender
		case 4:
			message.Method = "v1/" + strconv.Itoa(rand.Intn(99)) + "/ping"
		case 5:
			nonce = uuid.New().String()
			message.ID = []byte(nonce)
			message.Nonce = nonce
		}
		key_receiver = generateKey()
		crypt.Key = key_receiver
		//crypt.KeyUnlockDefault()
		ok := message.INVerifyMessage()
		if ok == true {
			t.Error("message.INVerifyMessage(): true - but content changed!")
			t.Fail()
		}
	}
}
