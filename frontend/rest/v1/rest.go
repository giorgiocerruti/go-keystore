package v1

import (
	"net/http"

	"github.com/giorgiocerruti/go-keystore/core"
	"github.com/gorilla/mux"
)

type restDFrontend struct {
	store *core.KeyValueStore
}

func (f *restDFrontend) Start(store *core.KeyValueStore, listen string) error {
	r := mux.NewRouter()
	//Register KeyValueHadler
	//matching "/v1/{key}"
	r.HandleFunc("/v1/{key}", f.KeyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", f.KeyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", f.KeyValueDeleteHandler).Methods("DELETE")

	return http.ListenAndServe(listen, r)
}
