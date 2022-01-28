package methods

import "x.x/x/deweb/lib"

func HandlePingV1(x lib.TransportStruct) {
	// pong is response, in this case we don't send anything back.
	if x.Data == "pong" {
		return
	}
}
