//go:build !badger
// +build !badger

package justdb

import (
	"crypto/sha1"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type MultiString []string

func (s *MultiString) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return errors.New("failed to scan multistring field - source is not a string")
	}
	*s = strings.Split(str, ",")
	return nil
}

func (s MultiString) Value() (driver.Value, error) {
	if s == nil || len(s) == 0 {
		return nil, nil
	}
	return strings.Join(s, ","), nil
}

type MapStringString map[string]string

func (s *MapStringString) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return errors.New("failed to scan multistring field - source is not a string")
	}
	return json.Unmarshal([]byte(str), s)
}

func (s MapStringString) Value() (driver.Value, error) {
	if s == nil || len(s) == 0 {
		return nil, nil
	}
	b, err := json.Marshal(s)
	return string(b), err
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
