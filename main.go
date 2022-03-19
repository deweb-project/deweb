package main

import (
	"runtime"

	"x.x/x/deweb/crypt"
	"x.x/x/deweb/frontend"
	"x.x/x/deweb/gui"
	"x.x/x/deweb/justdb"
	"x.x/x/deweb/lib"
	"x.x/x/deweb/server"
)

func main() {
	if runtime.GOOS == "linux" {
		justdb.Setup("/dev/shm/tmp")
	} else {
		panic("Non-Linux not implemented. // need correct path for DB")
	}
	crypt.LoadSelfKey()
	lib.LoadQueue()
	lib.LoadConversationList()

	server := server.NewServerLocal(50000)
	lib.SelfIdentifier = server.DEID.Identifier
	lib.SelfProto = server.DEID.Protocol
	go func() { panic(server.Start()) }()
	go frontend.Load()

	gui.Load()
}
