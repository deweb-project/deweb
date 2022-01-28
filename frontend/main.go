package frontend

import (
	"embed"
	"log"
	"net/http"
)

//go:embed *.css *.html *.js
var files embed.FS

func Load() {
	http.Handle("/", http.FileServer(http.FS(files)))

	// listen to port
	log.Println("Listening on :5313")
	http.ListenAndServe(":5313", nil)
}
