package main

import (
	"gorouter"
	"log"
	"net/http"
)

func main() {
	mux := gorouter.New()
	mux.GET("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(([]byte("Hello world")))
	})
	log.Fatal(http.ListenAndServe(":8000", mux))
}
