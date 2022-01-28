package lib

import "x.x/x/deweb/crypt"

type SelfID struct {
	OK bool
	ID string
}

func GetSelfID() SelfID {
	return SelfID{
		OK: true,
		ID: "dummyproto:aaaaa-aaaaa-aaaaaa-aaaaa-aaaaa[key=" + crypt.Key.GetFingerprint() + "]",
	}
}
