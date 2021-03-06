package v1

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/giorgiocerruti/go-keystore/core"
	"github.com/gorilla/mux"
)

func (f *RestFrontend) KeyValuePutHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["key"]

	//Read the body as it's a Reader interface
	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	fmt.Println(string(value), key)
	err = f.store.Put(key, string(value))
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
	}
	f.store.Put(key, string(value))
	w.WriteHeader(http.StatusCreated)
}

func (f *RestFrontend) KeyValueGetHandler(w http.ResponseWriter, r *http.Request) {

	var statusCode int
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := f.store.Get(key)
	if err != nil {
		if errors.Is(err, core.ErrorNoSuchKey) {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusInternalServerError
		}

		http.Error(w,
			err.Error(),
			statusCode)
	}

	w.Write([]byte(value))
}

func (f *RestFrontend) KeyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["key"]

	err := f.store.Delete(key)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
	}
	f.store.Delete(key)
	w.WriteHeader(http.StatusOK)

}
