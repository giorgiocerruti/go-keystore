package main

import (
	"log"
	"net/http"

	"github.com/giorgiocerruti/go-keystore/pkg/server"
)

const LISTEN_ADDRESS = ":8080"

func main() {
	r := server.NewRouter()

	log.Fatal(http.ListenAndServe(LISTEN_ADDRESS, r))
}
