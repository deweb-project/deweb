package lib

import (
	"math/rand"

	"x.x/x/deweb/justdb"
)

type QueueStore struct {
	ID    []byte
	Tasks map[string][]string
}

var Queue = QueueStore{
	ID: []byte("account1"),
}

func LoadQueue() {
	justdb.Read(&Queue)
}

func QueueTask(x TransportStruct) {
	justdb.Write(&x)
	Queue.Tasks[x.Destination] = append(Queue.Tasks[x.Destination], string(x.ID))
	justdb.Write(&Queue)
}

func GetTask(deid string) (x TransportStruct) {
	id := Queue.Tasks[deid][rand.Intn(len(Queue.Tasks[deid]))]
	x.ID = []byte(id)
	justdb.Read(&x)
	x.Tries++
	justdb.Write(&x)
	return
}

func RemoveTask(deid string, uuid string) {
	id := -1
	for i, v := range Queue.Tasks[deid] {
		if v == uuid {
			id = i
			break
		}

	}
	if id == -1 {
		return
	}
	Queue.Tasks[deid] = removeString(Queue.Tasks[deid], id)
	var x TransportStruct
	x.ID = []byte(uuid)
	justdb.Read(&x)
	justdb.Delete(&x)
	justdb.Write(&Queue)
}
func removeString(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
