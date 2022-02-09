package crypt

import (
	"log"
	"math/rand"
	"time"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
	"x.x/x/deweb/justdb"
)

type SelfKeyStore struct {
	ID   []byte
	Seed string
	Key  string
}

var Key *crypto.Key

func GetKey() *crypto.Key {
	Key, err := Key.Unlock([]byte{})
	if err != nil {
		panic(err)
	}
	return Key
}

var SelfKey = SelfKeyStore{
	ID: []byte("selfkey"),
}

var rsaBits = 2048

func LoadSelfKey() {
	justdb.Read(&SelfKey)

	InitPseudoRand()
	if SelfKey.Seed != "" {
		var err error
		Key, err = crypto.NewKeyFromArmored(SelfKey.Key)
		if err != nil {
			panic(err)
		}
		//KeyUnlockDefault()
		return
	}
	SelfKey.Seed = "not implemented"
	var err error
	print("Generating key.... ")
	SelfKey.Key, err = helper.GenerateKey("na", "na@na.na", []byte{}, "rsa", rsaBits)
	if err != nil {
		panic(err)
	}
	Key, err = crypto.NewKeyFromArmored(SelfKey.Key)
	//SelfKey.Key, err = Key.Armor()
	if err != nil {
		panic(err)
	}
	justdb.Write(&SelfKey)
	print("OK\n")
	//KeyUnlockDefault()
	//print("\n\nyour seed: '" + SelfKey.Seed + "'\n\n")
	//print(SelfKey.Key.CanEncrypt(), "\n")
}

func InitPseudoRand() {
	log.Println("init(rand): looking for correct seed")
	var tries int64 = 1024 * 16
	for {
		seed := time.Now().UnixNano() / tries
		rand.Seed(seed)
		if seed%tries == 0 {
			log.Println("init(rand): found correct seed", seed, tries)
			return
		}
		tries++
		if tries > 1024*64 {
			log.Println("init(rand): Max hit, resetting to 1024*16")
			tries = 1024 * 16
		}
	}
}
