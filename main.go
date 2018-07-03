package main

import (
	"encoding/json"
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
	r.HandleFunc("/hosts", makeHandler(listHosts)).Methods("GET")
	r.HandleFunc("/hosts/create", makeHandler(createHost)).Methods("POST")
	r.HandleFunc("/students", makeHandler(listStudents)).Methods("GET")
	r.HandleFunc("/students/create", makeHandler(createStudent)).Methods("POST")

	// App Engine routes incoming requests to the appropriate module on port 8080.
	// https://cloud.google.com/appengine/docs/flexible/custom-runtimes/build#listen_to_port_8080
	//
	http.ListenAndServe(":8080", r)
}

func listHosts(w http.ResponseWriter, req *http.Request) {
	rows, err := db.Get("first_name, last_name", "hosts")
	if err != nil {
		internalError(w, err.Error())
	}

	j, err := json.Marshal(rows)
	if err != nil {
		internalError(w, err.Error())
	}

	w.Write(j)
}

func createHost(w http.ResponseWriter, req *http.Request) {
	host := models.Host{req.PostForm}
	saved := save(host, w)
	if saved {
		log.Println("successfully inserted host: %v", host)
	}
}

func listStudents(w http.ResponseWriter, req *http.Request) {
	rows, err := db.Get("first_name, last_name, country_of_origin", "students")
	if err != nil {
		internalError(w, err.Error())
	}

	j, err := json.Marshal(rows)
	if err != nil {
		internalError(w, err.Error())
	}

	w.Write(j)
}

func createStudent(w http.ResponseWriter, req *http.Request) {
	student := models.Student{req.PostForm}
	saved := save(student, w)
	if saved {
		log.Println("successfully inserted student: %v", student)
	}
}

func makeHandler(fn func(w http.ResponseWriter, req *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			internalError(w, err.Error())
			return
		}

		fn(w, req)
	}
}

func save(model models.Model, w http.ResponseWriter) bool {
	validationErrors, databaseError := db.Insert(model)
	if databaseError != nil {
		internalError(w, databaseError.Error())
		return false
	} else if !validationErrors.None() {
		badRequest(w, validationErrors.Concatenate(", "))
		return false
	}

	return true
}
