package server

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	//Register KeyValueHadler
	//matching "/v1/{key}"
	r.HandleFunc("/v1/{key}", KeyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", KeyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", KeyValueDeleteHandler).Methods("DELETE")

	return r
}
