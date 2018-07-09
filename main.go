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
	r.HandleFunc("/hosts", indexHandler(&models.Host{})).Methods("GET")
	r.HandleFunc("/hosts/create", insertHandler(&models.Host{})).Methods("POST")
	r.HandleFunc("/students", indexHandler(&models.Student{})).Methods("GET")
	r.HandleFunc("/students/create", insertHandler(&models.Student{})).Methods("POST")

	// App Engine routes incoming requests to the appropriate module on port 8080.
	// https://cloud.google.com/appengine/docs/flexible/custom-runtimes/build#listen_to_port_8080
	//
	http.ListenAndServe(":8080", r)
}

func indexHandler(model models.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		rows, err := db.Get(model)
		if err != nil {
			internalError(w, err.Error())
		}

		j, err := json.Marshal(rows)
		if err != nil {
			internalError(w, err.Error())
		}

		w.Write(j)
	}
}

func insertHandler(model models.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			internalError(w, err.Error())
			return
		}

		model.Init(req.PostForm)

		validationErrors, databaseError := db.Insert(model)
		if databaseError != nil {
			internalError(w, databaseError.Error())
			return
		} else if !validationErrors.None() {
			badRequest(w, validationErrors.Concatenate(", "))
			return
		}

		log.Println("successfully inserted record: %v", model)
	}
}
