package v1

import (
	"net/http"
	"os"

	"github.com/giorgiocerruti/go-keystore/core"
	"github.com/gorilla/mux"
)

type RestFrontend struct {
	store *core.KeyValueStore
}

func (f *RestFrontend) Start(store *core.KeyValueStore) error {
	r := mux.NewRouter()
	f.store = store

	//Register KeyValueHadler
	//matching "/v1/{key}"
	r.HandleFunc("/v1/{key}", f.KeyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", f.KeyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", f.KeyValueDeleteHandler).Methods("DELETE")

	return http.ListenAndServe(os.Getenv("TLOG_REST_ADDR"), r)
}

func NewRestFrontend() *RestFrontend {
	return &RestFrontend{}
}
