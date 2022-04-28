package main

import (
	"log"
	"os"

	"github.com/giorgiocerruti/go-keystore/core"
	"github.com/giorgiocerruti/go-keystore/frontend"
	"github.com/giorgiocerruti/go-keystore/transact"
)

func main() {
	//Create our TransactionLogger
	tl, err := transact.NewTransactionLogger(os.Getenv("TLOG_TYPE"))
	if err != nil {
		log.Fatal(err)
	}

	//Creare the core and tell it wich TL to use
	store := core.NewKeyValueStore(tl)
	err = store.Restore()
	if err != nil {
		log.Fatal(err)
	}

	//Create the frontend
	fe, err := frontend.NewFrontEnd(os.Getenv("TLOG_FRONTEND"))
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(fe.Start(store))

}
