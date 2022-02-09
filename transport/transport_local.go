package transport

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"strings"

	"x.x/x/deweb/crypt"
	"x.x/x/deweb/lib"
)

var ErrLocalNotLocal = errors.New("protocol is not local")

var connectionLocal = Connection{
	Protocol:             "local",
	Destination:          "",
	Key:                  "",
	EstabilishConnection: func() error { return nil }, // no need to estabilish connection, at least now.
	Send: func(deid lib.DEID, key string, x lib.TransportStruct) error {
		if deid.Protocol != "local" {
			return ErrLocalNotLocal
		}
		addr, err := net.ResolveTCPAddr("tcp", deid.Identifier)
		if err != nil {
			return err
		}
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			return err
		}
		x.OUTAttachPublicKey()
		x.OUTAttachSignature()
		log.Println("Sending:", x.String())
		log.Println(key)
		data := crypt.Encrypt(x, key)
		conn.Write(data)
		Handleconn(conn, false)
		return nil
	},
	GetKey: func(deid lib.DEID) (key string, err error) {
		addr, err := net.ResolveTCPAddr("tcp", deid.Identifier)
		if err != nil {
			panic(err)
		}
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		_, err = conn.Write([]byte("plaintext/init/v1\n"))
		if err != nil {
			panic(err)
		}
		all, err := ioutil.ReadAll(conn)
		if err != nil {
			panic(err)
		}
		// all = string
		// 1. lib.GetSelfID()
		// 2-.... crypt.GetKey
		s := strings.SplitN(string(all), "\n", 2)
		if len(s) != 2 {
			log.Println(s[0], len(s))
			panic(s)
		}
		if s[0] != deid.String() {
			log.Println(s[0], deid.String())
			log.Println(errors.New("(soon)panic: incorrect response. deid missmatch"))
		}
		return s[1], nil
	},
}

func GetTransportLocal(deid lib.DEID) Connection {
	conn_local := connectionLocal
	conn_local.Destination = deid.String()
	key, err := conn_local.GetKey(deid)
	if err != nil {
		panic(err)
	}
	conn_local.Key = key

	log.Println("received key")
	return conn_local
}

//
