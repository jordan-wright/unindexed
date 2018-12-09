package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jordan-wright/unindexed"
)

func main() {
	router := mux.NewRouter()
	router.PathPrefix("/").Handler(http.FileServer(unindexed.Dir("../static")))
	log.Fatal(http.ListenAndServe(":8080", router))
}
