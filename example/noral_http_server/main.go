package main

import (
	"log"
	"net/http"
	// "github.com/xujiajun/gorouter"
)

func IndexHandleFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.Path))
}

func main() {
	http.HandleFunc("/", IndexHandleFunc)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
