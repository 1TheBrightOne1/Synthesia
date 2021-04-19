package server

import (
	"log"
	"net/http"
)

func Start() {
	fs := http.FileServer(http.Dir("./server/static"))
	sessions := http.FileServer(http.Dir("/outputs/sessions"))
	http.HandleFunc("/download/", download)
	http.HandleFunc("/setup/", setup)
	http.Handle("/files/", http.StripPrefix("/files/", fs))
	http.Handle("/sessions/", http.StripPrefix("/sessions/", sessions))
	http.HandleFunc("/generate/", generate)
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe("0.0.0.0:8083", nil))
}
