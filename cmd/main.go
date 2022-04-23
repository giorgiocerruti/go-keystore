package main

import (
	"fmt"
	"log"
	"net/http"

	store "github.com/giorgiocerruti/go-keystore/pkg/api/v1"
	"github.com/giorgiocerruti/go-keystore/pkg/server"
)

const LISTEN_ADDRESS = ":8080"

func main() {
	fmt.Println("Intializing...")

	tLogger, err := store.InitializeTransdactionLogger()
	if err != nil {
		fmt.Printf("error initializing the storage: %s", err)
	}

	r := server.NewRouter(tLogger)
	fmt.Printf("Server listening %s", LISTEN_ADDRESS)

	log.Fatal(http.ListenAndServe(LISTEN_ADDRESS, r))

}
