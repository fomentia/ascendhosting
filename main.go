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
	r.HandleFunc("/hosts/create", createHost).Methods("POST")
	r.HandleFunc("/students/create", createStudent).Methods("POST")

	// App Engine routes incoming requests to the appropriate module on port 8080.
	// https://cloud.google.com/appengine/docs/flexible/custom-runtimes/build#listen_to_port_8080
	//
	http.ListenAndServe(":8080", r)
}

func createHost(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		internalError(w, err.Error())
		return
	}

	firstName := req.PostForm.Get("firstName")
	lastName := req.PostForm.Get("lastName")
	host := models.Host{firstName, lastName}

	validationErrors, databaseError := db.Insert(host)
	if databaseError != nil {
		internalError(w, err.Error())
		return
	} else if !validationErrors.None() {
		badRequest(w, validationErrors.Concatenate(", "))
		return
	}

	log.Println("successfully inserted host: %v", host)
}

func createStudent(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		internalError(w, err.Error())
		return
	}

	firstName := req.PostForm.Get("firstName")
	lastName := req.PostForm.Get("lastName")
	country := req.PostForm.Get("countryOfOrigin")
	student := models.Student{firstName, lastName, country}

	validationErrors, databaseError := db.Insert(student)
	if databaseError != nil {
		internalError(w, err.Error())
		return
	} else if !validationErrors.None() {
		badRequest(w, validationErrors.Concatenate(", "))
		return
	}

	log.Println("successfully inserted student: %v", student)
}
