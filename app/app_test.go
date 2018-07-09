package app

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fomentia/ascendhosting/database"
	"github.com/fomentia/ascendhosting/models"
)

type mockDB struct {
	Rows []database.Row
}

func (db *mockDB) Insert(model models.Model) (validationErrors models.Errors, databaseError error) {
	log.Println("Insert called on mockDB")
	return models.Errors{}, nil
}

func (db *mockDB) Get(model models.Model) ([]database.Row, error) {
	return []database.Row{
		database.Row{
			"firstName": "Peter",
			"lastName":  "Parker",
		},
	}, nil
}

func TestGetHostsHandler(t *testing.T) {
	req, reqErr := http.NewRequest("GET", "/hosts", nil)
	if reqErr != nil {
		t.Fatal(reqErr)
	}

	rr := httptest.NewRecorder()
	DB = &mockDB{}
	handler := IndexHandler(&models.Host{})
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Wanted response code %v, but got %v", http.StatusOK, rr.Code)
	}
}
