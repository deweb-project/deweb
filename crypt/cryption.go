package crypt

import (
	"encoding/json"
	"log"

	"github.com/ProtonMail/gopenpgp/v2/helper"
)

// This function receive an object and json it.
func Encrypt(x interface{}, key_armored string) []byte {
	rawjson, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}

	encrypted, err := helper.EncryptBinaryMessageArmored(key_armored, rawjson)
	if err != nil {
		panic(err)
	}
	return []byte(encrypted)
}
func Decrypt(encrypted []byte, x interface{}) error {
	privkey_armored, err := Key.Armor()
	if err != nil {
		return err
	}
	rawjson, err := helper.DecryptBinaryMessageArmored(privkey_armored, []byte{}, string(encrypted))
	if err != nil {
		return err
	}
	err = json.Unmarshal(rawjson, x)
	if err != nil {
		return err
	}
	log.Println("Decrypt:", string(rawjson))
	return nil
}

func init() {
	log.SetFlags(log.Lshortfile)
}
