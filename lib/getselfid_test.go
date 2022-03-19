package lib_test

import (
	"testing"

	"x.x/x/deweb/crypt"
	"x.x/x/deweb/justdb"
	"x.x/x/deweb/lib"
)

func TestGetSelfID(t *testing.T) {
	justdb.Setup("/dev/shm/tmp")
	crypt.LoadSelfKey()
	id := lib.GetSelfID()
	if !id.OK {
		t.Fail()
	}
}
