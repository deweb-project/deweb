//go:build badger
// +build badger

package justdb

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	badger "github.com/dgraph-io/badger/v3"
)

var DB *badger.DB
var ErrKeyNotFound = badger.ErrKeyNotFound

func Setup(path string) {

	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		panic(err)
	}
	DB = db
again:
	err = DB.RunValueLogGC(0.7)
	if err == nil {
		goto again
	}
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
		again:
			err := DB.RunValueLogGC(0.7)
			if err == nil {
				goto again
			}
		}
	}()
}

func Close() error {
	return DB.Close()
}

type DataModel struct {
	ID []byte
}

func Hash(data []byte) []byte {
	hasher := sha1.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}

func Delete(data interface{}) {
	//log.Printf("%T", data)
	datatype := Hash([]byte(fmt.Sprintf("%T", data)))[0:8]
	id := reflect.ValueOf(data).Elem().FieldByName("ID").Interface().([]byte)
	if len(id) < 20 {
		id = Hash(id)
	}
	txn := DB.NewTransaction(true)
	err := txn.Delete(append(datatype, id...))
	if err != nil {
		panic(err)
	}
	err = txn.Commit()
	if err != nil {
		panic(err)
	}
}

func Write(data interface{}) {
	//log.Printf("%T", data)
	datatype := Hash([]byte(fmt.Sprintf("%T", data)))[0:8]
	id := reflect.ValueOf(data).Elem().FieldByName("ID").Interface().([]byte)
	if len(id) < 20 {
		id = Hash(id)
	}
	writedata, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	txn := DB.NewTransaction(true)
	//log.Println(append(datatype, id...))
	err = txn.Set(append(datatype, id...), writedata)
	if err != nil {
		panic(err)
	}
	err = txn.Commit()
	if err != nil {
		panic(err)
	}
}

func Read(data interface{}) error {
	//log.Printf("%T", data)
	datatype := Hash([]byte(fmt.Sprintf("%T", data)))[0:8]
	id := reflect.ValueOf(data).Elem().FieldByName("ID").Interface().([]byte)
	if len(id) < 20 {
		id = Hash(id)
	}
	txn := DB.NewTransaction(false)
	//log.Println(append(datatype, id...))
	item, err := txn.Get(append(datatype, id...))
	if err != nil {
		return err
	}
	return item.Value(func(val []byte) error {
		return json.Unmarshal(val, data)
	})
}
