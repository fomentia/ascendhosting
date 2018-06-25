package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var db *DB

func main() {
	var err error
	db, err = initDB()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/hosts/create", root).Methods("POST")

	http.ListenAndServe(":5000", r)
}

func root(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		internalError(w, err.Error())
		return
	}

	firstName := req.PostForm.Get("first_name")
	if len(firstName) == 0 {
		badRequest(w, "first_name is blank")
		return
	}

	lastName := req.PostForm.Get("last_name")
	if len(lastName) == 0 {
		badRequest(w, "last_name is blank")
		return
	}

	host := Host{firstName, lastName}
	err = db.insertHost(host)
	if err != nil {
		internalError(w, err.Error())
	}

	log.Println("successfully inserted host: %v", host)
}

func internalError(w http.ResponseWriter, message string) {
	log.Println(message)
	http.Error(w, message, http.StatusInternalServerError)
}

func badRequest(w http.ResponseWriter, message string) {
	log.Println(message)
	http.Error(w, message, http.StatusBadRequest)
}
