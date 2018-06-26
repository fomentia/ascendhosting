package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/fomentia/ascendhosting/database"
	"github.com/fomentia/ascendhosting/models"
)

var db *database.DB

func main() {
	var err error
	db, err = database.InitDB()
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
	host := models.Host{firstName, lastName}

	errors := db.Insert(host)
	if !errors.None() {
		badRequest(w, errors.Concatenate(", "))
		return
	}

	log.Println("successfully inserted host: %v", host)
}
