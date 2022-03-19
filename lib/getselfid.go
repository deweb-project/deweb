package lib

import (
	"x.x/x/deweb/crypt"
)

type SelfID struct {
	OK bool
	ID string
}

var SelfIdentifier = ""
var SelfProto = ""

func GetSelfID() SelfID {
	return SelfID{
		OK: true,
		ID: SelfProto + ":" + SelfIdentifier + "[key=" + crypt.Key.GetFingerprint() + "]",
	}
}
