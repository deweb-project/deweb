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
	ui.Bind("createConversation", lib.CreateConversation)
	ui.Bind("ignoreInvitation", lib.IgnoreInvitation)
	ui.Bind("announceConversation", lib.AnnounceConversation)
	ui.Bind("banUser", lib.BanUser)
	ui.Bind("getUser", lib.GetUser)
	ui.Bind("loadConversationList", lib.LoadConversationList)
	ui.Bind("queueTask", lib.QueueTask)
	ui.Bind("sendInvitation", lib.SendInvitation)
	ui.Bind("sendMessage", lib.SendMessage)
	ui.Bind("acceptInvitation", lib.AcceptInvitation)

	ui.Load("http://localhost:5313")
	<-ui.Done()
}
