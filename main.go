package main

import (
	"runtime"

	"x.x/x/deweb/crypt"
	"x.x/x/deweb/frontend"
	"x.x/x/deweb/gui"
	"x.x/x/deweb/justdb"
)

func main() {
	if runtime.GOOS == "linux" {
		justdb.Setup("/dev/shm/tmp")
	} else {
		justdb.Setup("")
	}
	crypt.LoadSelfKey()
	go frontend.Load()
	gui.Load()
}
