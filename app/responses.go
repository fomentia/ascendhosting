package app

import (
	"log"
	"net/http"
)

func responseError(w http.ResponseWriter, err error, code int) {
	log.Println(err.Error())
	http.Error(w, err.Error(), code)
}
