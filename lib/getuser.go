package lib

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"

	"x.x/x/deweb/justdb"
)

type UserConnectionStatus struct {
	CanConnect       bool // do we know how?
	LastSeen         int
	AlreadyConnected bool
}
type UserInfo struct {
	ID           []byte // string(deid)
	ConnectionID int
	Connection   UserConnectionStatus
	DEIDID       int
	DEID         DEID
}

type DEID struct {
	Protocol   string
	Identifier string
	Key        string
	Extra      map[string]string
}

var ParseDEIDIncorrectAddress = "parsedeid: incorrect address given"

func ParseDEID(deid string) (DEID, error) {
	deid_split := strings.SplitN(deid, ":", 2)
	log.Println(deid_split)
	if len(deid_split) != 2 {
		print("1. len(deid_split) != 2\n")
		return DEID{}, errors.New(ParseDEIDIncorrectAddress)
	}
	protocol := deid_split[0]
	deid_split = strings.SplitN(deid_split[1], "[", 2)
	if len(deid_split) != 2 {
		print("2. len(deid_split) != 2\n")
		return DEID{}, errors.New(ParseDEIDIncorrectAddress)
	}
	identifier := deid_split[0]
	if len(deid_split[1]) == 0 {
		print("3. len(deid_split[1]) == 0\n")
		return DEID{}, errors.New(ParseDEIDIncorrectAddress)
	}
	vals, err := url.ParseQuery(deid_split[1][0 : len(deid_split[1])-1])
	if err != nil {
		return DEID{}, err
	}
	if !vals.Has("key") {
		print("4. !vals.Has(\"key\")\n")
		return DEID{}, errors.New(ParseDEIDIncorrectAddress)
	}
	var finalvals = make(map[string]string)
	for k := range vals {
		if k == "key" {
			continue
		}
		finalvals[k] = vals.Get(k)
	}

	return DEID{
		Protocol:   protocol,
		Identifier: identifier,
		Key:        vals.Get("key"),
		Extra:      finalvals,
	}, nil
}

type GetUserResp struct {
	OK       bool
	Error    MyError
	UserInfo UserInfo
}

func GetUser(deid string) GetUserResp {

	// id format:
	// proto:identifier[key=asdasdasdasdasdasdasd]
	// Example:
	// libp2p:/ip4/7.7.7.7/tcp/4242/p2p/QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N[key=???]
	// tor:asdasd.onion[key=???]
	// proxied:/ip4/10.8.42.42/tcp/4242/p2p/QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N[key=???&via=tor:asdasd.onion[key=???]]
	// Keep in mind that deweb assumes plaintext connection and encrypts all data anyway.
	//So even if we connect over encrypted tunel using libp2p/tor/i2p data is still being
	//encrypted and decrypted by the deweb daemon, just as the connection was plaintext.
	//This is done to ensure that even if the the connection method becomes compromised,
	//you won't lose encryption.
	// Yet you don't need to trust proxy. All it receive is destination id and encrypted packet.
	var err error
	var user = UserInfo{
		ID: []byte(deid),
	}
	justdb.Read(&user)
	user.DEID, err = ParseDEID(deid)
	if err != nil {
		justdb.Write(&user)
	}
	return GetUserResp{
		OK:       err == nil,
		Error:    MyError{err},
		UserInfo: user,
	}
}

type MyError struct {
	error
}

func (me MyError) MarshalJSON() ([]byte, error) {
	return json.Marshal(me.Error())
}
