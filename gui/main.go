package gui

import (
	"github.com/zserge/lorca"
	"x.x/x/deweb/lib"
)

func Load() {
	ui, err := lorca.New("", "", 1280, 720)
	if err != nil {
		panic(err)
	}
	// To avoid binding too many functions, that aren't really used, and could
	//cause security issues, please include functions as needed, following
	//this schema:
	ui.Bind("getSelfID", lib.GetSelfID)
	ui.Load("http://localhost:5313")
	<-ui.Done()
}
