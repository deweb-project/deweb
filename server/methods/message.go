package methods

import (
	"x.x/x/deweb/lib"
)

func HandleMessageV1(x lib.TransportStruct) {
	lib.HandleNewMessage(x)
}
