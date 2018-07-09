package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/fomentia/ascendhosting/database"
	"github.com/fomentia/ascendhosting/models"
)

var DB database.Database

func Serve(database database.Database) {
	DB = database

	r := mux.NewRouter()
	r.HandleFunc("/hosts", IndexHandler(&models.Host{})).Methods("GET")
	r.HandleFunc("/hosts/create", InsertHandler(&models.Host{})).Methods("POST")
	r.HandleFunc("/students", IndexHandler(&models.Student{})).Methods("GET")
	r.HandleFunc("/students/create", InsertHandler(&models.Student{})).Methods("POST")

	// App Engine routes incoming requests to the appropriate module on port 8080.
	// https://cloud.google.com/appengine/docs/flexible/custom-runtimes/build#listen_to_port_8080
	//
	http.ListenAndServe(":8080", r)
}

func IndexHandler(model models.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		rows, err := DB.Get(model)
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

func InsertHandler(model models.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			internalError(w, err.Error())
			return
		}

		model.Init(req.PostForm)

		validationErrors, databaseError := DB.Insert(model)
		if databaseError != nil {
			internalError(w, databaseError.Error())
			return
		} else if !validationErrors.None() {
			badRequest(w, validationErrors.Concatenate(", "))
			return
		}

		log.Printf("successfully inserted record: %v\n", model)
	}
}
