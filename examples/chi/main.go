package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jordan-wright/unindexed"
)

func main() {
	router := chi.NewRouter()
	router.Handle("/*", http.FileServer(unindexed.Dir("../static/")))
	log.Fatal(http.ListenAndServe(":8080", router))
}
