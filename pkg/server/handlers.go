package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	gapi "github.com/giorgiocerruti/go-keystore/pkg/api"

	"github.com/gorilla/mux"
)

func KeyValuePutHandler(w http.ResponseWriter, r *http.Request) {
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
	err = gapi.Put(key, string(value))
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

func KeyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := gapi.Get(key)
	if err != nil {
		if errors.Is(err, gapi.ErrorNoSuchKey) {
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
