package lib_test

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"x.x/x/deweb/crypt"
	"x.x/x/deweb/justdb"
	"x.x/x/deweb/lib"
)

func jsonprint(x interface{}) string {
	b, _ := json.MarshalIndent(x, "", "    ")
	return string(b)
}

// SECURITY:
// Chances are that this function is fundamentally wrong.
// Something in here is just not right.
func TestQueue(t *testing.T) {
	err := os.RemoveAll("/dev/shm/tmp_queue")
	if err != nil {
		log.Println("WARN: Failed to delete /dev/shm/tmp_queue")
	}
	justdb.Setup("/dev/shm/tmp_queue")
	crypt.LoadSelfKey()
	lib.LoadQueue()
	var nonce = uuid.New().String()
	var sometask = lib.TransportStruct{
		ID:          []byte(nonce),
		Source:      lib.GetSelfID().ID,
		Destination: lib.GetSelfID().ID,
		Nonce:       nonce,
		Method:      "v1/0/ping",
	}
	lib.QueueTask(sometask)
	task := lib.GetTask(lib.GetSelfID().ID) // first task is something else?
	task = lib.GetTask(lib.GetSelfID().ID)
	task.Tries--
	if !reflect.DeepEqual(sometask, task) {
		log.Println("Not equal!")
		log.Println("sometask:", jsonprint(sometask))
		log.Println("task:", jsonprint(task))
		t.Fail()
	}
	lib.RemoveTask(lib.GetSelfID().ID, task.Nonce)
	task = lib.GetTask(lib.GetSelfID().ID)
	if reflect.DeepEqual(sometask, task) {
		log.Println("Equal! (But it shouldn't be)")
		t.Fail()
	}
}
