package main

import (
	"log"
	"net/http"
)

func internalError(w http.ResponseWriter, message string) {
	log.Println(message)
	http.Error(w, message, http.StatusInternalServerError)
}

func badRequest(w http.ResponseWriter, message string) {
	log.Println(message)
	http.Error(w, message, http.StatusBadRequest)
}
