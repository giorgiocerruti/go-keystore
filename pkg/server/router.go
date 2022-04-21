package server

import (
	gapi "github.com/giorgiocerruti/go-keystore/api"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	//Register KeyValueHadler
	//matching "/v1/{key}"
	r.HandleFunc("/v1/{key}", gapi.KeyValuePutHandler).Method("PUT")

	return r
}
