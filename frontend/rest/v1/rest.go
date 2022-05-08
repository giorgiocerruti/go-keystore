package v1

import (
	"fmt"
	"net/http"

	"github.com/giorgiocerruti/go-keystore/core"
	"github.com/gorilla/mux"
)

type RestFrontend struct {
	store  *core.KeyValueStore
	Config RestConfig
}

type RestConfig struct {
	Address string
	Port    string
}

func (f *RestFrontend) Start(store *core.KeyValueStore) error {
	r := mux.NewRouter()
	f.store = store
	if f.Config.Port == "" {
		f.Config.Port = "8080"
	}

	addr := fmt.Sprintf("%s:%s", f.Config.Address, f.Config.Port)

	//Register KeyValueHadler
	//matching "/v1/{key}"
	r.HandleFunc("/v1/{key}", f.KeyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", f.KeyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", f.KeyValueDeleteHandler).Methods("DELETE")

	fmt.Printf("Listening on address:port %s", addr)
	return http.ListenAndServe(addr, r)
}

func NewRestFrontend() *RestFrontend {
	return &RestFrontend{}
}
