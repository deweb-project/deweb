package lib

import (
	"strings"

	"x.x/x/deweb/justdb"
)

type QueueStore struct {
	ID    []byte
	Tasks justdb.MultiString `gorm:"type:text"`
}

var Queue = QueueStore{
	ID: []byte("account1"),
}

func LoadQueue() {
	justdb.Read(&Queue)
}

func QueueTask(x TransportStruct) {
	x.OUTInitNonce()
	x.OUTAttachPublicKey()
	x.OUTAttachSignature()
	justdb.Write(&x)
	Queue.Tasks = append(Queue.Tasks, string(x.Destination)+"|||"+string(x.ID))
	justdb.Write(&Queue)
}

func GetTask(deid string) (x TransportStruct) {
	for _, v := range Queue.Tasks {
		s := strings.Split(v, "|||")
		if s[0] == deid {
			id := s[1]
			x.ID = []byte(id)
			justdb.Read(&x)
			x.Tries++
			justdb.Write(&x)
			return
		}
	}
	return x
}

func RemoveTask(deid string, uuid string) {
	id := -1

	for i, v := range Queue.Tasks {
		s := strings.Split(v, "|||")
		if s[0] == deid && s[1] == uuid {
			id = i
			break
		}
	}
	if id == -1 {
		return
	}
	Queue.Tasks = removeString(Queue.Tasks, id)
	var x TransportStruct
	x.ID = []byte(uuid)
	justdb.Read(&x)
	justdb.Delete(&x)
	justdb.Write(&Queue)
}
func removeString(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
