//go:build !badger
// +build !badger

package justdb

import (
	"crypto/sha1"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MultiMultiByte [][]byte

func (s *MultiMultiByte) Scan(src interface{}) error {
	str, ok := src.([]byte)
	if !ok {
		log.Fatal("Failed to parse MultiMultiByte")
	}
	err := json.Unmarshal(str, s)
	if err != nil {
		log.Fatal("Failed to umarshal")
	}
	return nil
}

func (s MultiMultiByte) Value() (driver.Value, error) {
	return json.Marshal(s)
}

var DB *gorm.DB

type MultiString []string

func (s *MultiString) Scan(src interface{}) error {
	str, ok := src.([]byte)
	if !ok {
		log.Fatal("Failed to parse MultiMultiByte")
	}
	err := json.Unmarshal(str, s)
	if err != nil {
		log.Fatal("Failed to umarshal")
	}
	return nil
}

func (s MultiString) Value() (driver.Value, error) {
	return json.Marshal(s)
}

type MapStringString map[string]string

func (s *MapStringString) Scan(src interface{}) error {
	str, ok := src.([]byte)
	if !ok {
		log.Fatal("Failed to parse MultiMultiByte")
	}
	err := json.Unmarshal(str, s)
	if err != nil {
		log.Fatal("Failed to umarshal")
	}
	return nil
}

func (s MapStringString) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func Setup(path string) {
	var err error
	os.MkdirAll(path, 0750)
	DB, err = gorm.Open(sqlite.Open(path+"/sqlite.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}

func Close() error {
	return nil
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
	ivar(data)
	DB.Delete(data)
}

func Write(data interface{}) error {
	ivar(data)
	//id := reflect.ValueOf(data).Elem().FieldByName("ID").Interface().([]byte)
	//if len(id) < 20 {
	//	data = Hash(id)
	//}
	err := DB.Save(data).Error
	if err != nil {
		panic(err)
	}
	return err
}

func Read(data interface{}) error {
	ivar(data)
	//log.Printf("%T", data)
	id := reflect.ValueOf(data).Elem().FieldByName("ID").Interface().([]byte)
	//if len(id) < 20 {
	//	id = Hash(id)
	//}
	return DB.First(data, "ID = ?", id).Error
}

var inited = make(map[string]bool)

func ivar(data interface{}) {
	fmt.Printf("%T\n", data)
	datatype := string(Hash([]byte(fmt.Sprintf("%T", data)))[0:8])
	if inited[datatype] {
		return
	}
	inited[datatype] = true
	err := DB.AutoMigrate(data)
	if err != nil {
		panic(err)
	}
}
