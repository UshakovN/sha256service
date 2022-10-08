package tools

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, _ *http.Request, p any) {
	b, err := json.Marshal(p)
	if err != nil {
		log.Printf("cannot convert data to json: %v", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		log.Printf("cannot write to message writer: %v", err)
	}
}

func WriteRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("error in request %s: %v", r.RequestURI, err)
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func WriteInternalError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal error in request %s: %v", r.RequestURI, err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func HandleHealthRequest(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("/ok"))
	if err != nil {
		log.Printf("cannot write to response writer: %v", err)
	}
}
