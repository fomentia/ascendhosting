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

	firstName := req.PostForm.Get("firstName")
	lastName := req.PostForm.Get("lastName")
	host := Host{firstName, lastName}

	errors := db.insertHost(host)
	if !errors.none() {
		badRequest(w, errors.concatenate(", "))
		return
	}

	log.Println("successfully inserted host: %v", host)
}
