package justdb_test

import (
	"log"
	"os"
	"testing"

	"x.x/x/deweb/justdb"
	"x.x/x/deweb/lib"
)

type TypeByte struct {
	ID []byte
}

type TypeCommon struct {
	ID     []byte
	String string
	Int    int
	Uint   uint
	Bool   bool
}

type TypeSString struct {
	ID      []byte
	SString justdb.MultiString `gorm:"type:text"`
}

func TestTypes(t *testing.T) {
	err := os.RemoveAll("/dev/shm/testdb")
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	justdb.Setup("/dev/shm/testdb")
	err = justdb.Write(&TypeByte{
		ID: []byte("test"),
	})
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	err = justdb.Write(&TypeCommon{
		ID: []byte("test"),
	})
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	err = justdb.Write(TypeSString{
		ID:      []byte("test"),
		SString: []string{"String1", "String1", "String3"},
	})
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	// now let's test actual db schema
	err = justdb.Write(lib.UserInfo{
		ID: []byte("test"),
	})
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	err = justdb.Write(lib.QueueStore{
		ID: []byte("test"),
	})
	if err != nil {
		panic(err)
		//t.Fail()
	}
}
