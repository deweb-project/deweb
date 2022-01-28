package gui

import (
	"github.com/zserge/lorca"
	"x.x/x/deweb/lib"
)

func Load() {
	ui, _ := lorca.New("", "", 1280, 720)
	ui.Bind("getSelfID", lib.GetSelfID)
	ui.Bind("getUser", lib.GetUser)
	ui.Load("http://localhost:5313")
	<-ui.Done()
}
