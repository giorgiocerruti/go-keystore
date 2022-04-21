package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/giorgiocerruti/go-keystore/pkg/server"
)

const LISTEN_ADDRESS = ":8080"

func main() {
	r := server.NewRouter()

	fmt.Printf("Server listening %s", LISTEN_ADDRESS)
	log.Fatal(http.ListenAndServe(LISTEN_ADDRESS, r))
}
