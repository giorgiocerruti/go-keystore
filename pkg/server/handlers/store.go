package handlers

import (
	gapi "github/giorgiocerruti/go-keystore/api"
	"io"
	"net/http"

	"github.com/giorilla/mux"
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

	err = gapi.Put(key, string(value))
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}
