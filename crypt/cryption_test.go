package crypt_test

import (
	"log"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
	"x.x/x/deweb/crypt"
)

type TestInterface struct {
	String string
	Bool   bool
	Byte   byte
	SByte  []byte
	Uint64 uint64
	Int64  int64
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
func TestEncryptDecrypt(t *testing.T) {
	for v := int64(1); v < 100; v++ {
		rand.Seed(time.Now().UnixNano() / v)
		plain := TestInterface{
			String: RandomString(1024 * int(v)),
			Bool:   v%2 == 0,
			Byte:   byte(RandomString(1)[0]),
			SByte:  []byte(RandomString(1024 * int(v))),
			Uint64: uint64(time.Now().UnixNano() / v),
			Int64:  time.Now().UnixNano() / v,
		}

		// RSA, string
		rsaKey, err := helper.GenerateKey("na", "na@na.na", []byte{}, "rsa", 1024) // 1024 just for testing
		if err != nil {
			t.Error(err)
		}
		crypt.Key, err = crypto.NewKeyFromArmored(rsaKey)
		if err != nil {
			t.Error(err)
		}
		armored_receiver_public, err := crypt.Key.GetArmoredPublicKey()
		if err != nil {
			t.Error(err)
		}
		encrypted := crypt.Encrypt(plain, armored_receiver_public)

		//
		var decrypted TestInterface
		crypt.Decrypt(encrypted, &decrypted)
		log.Print(decrypted, "\n", plain)
		if !reflect.DeepEqual(plain, decrypted) {
			t.Fail()
		}
	}
}
