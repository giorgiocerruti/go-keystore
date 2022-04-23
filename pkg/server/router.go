package server

import (
	"github.com/giorgiocerruti/go-keystore/pkg/logger"
	"github.com/gorilla/mux"
)

func NewRouter(logger logger.TransactionLogger) *mux.Router {
	r := mux.NewRouter()
	//Register KeyValueHadler
	//matching "/v1/{key}"
	r.HandleFunc("/v1/{key}", KeyValuePutHandler(logger)).Methods("PUT")
	r.HandleFunc("/v1/{key}", KeyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", KeyValueDeleteHandler(logger)).Methods("DELETE")

	return r
}
