package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/fomentia/ascendhosting/database"
	"github.com/fomentia/ascendhosting/models"
)

type mockDB struct {
	Rows []database.Row
}

func (db *mockDB) Insert(model models.Model) (validationErrors models.Errors, databaseError error) {
	validationErrors = models.Validate(model)

	if validationErrors.None() {
		row := database.Row{}
		for index, column := range strings.Split(model.Columns(), ", ") {
			row[column] = model.Values()[index]
		}
		db.Rows = append(db.Rows, row)
	}

	return
}

func (db *mockDB) Get(model models.Model) ([]database.Row, error) {
	return db.Rows, nil
}

func TestHosts(t *testing.T) {
	DB = &mockDB{}

	requestParams := url.Values{"firstName": {"Peter"}, "lastName": {"Parker"}}
	req, reqErr := http.NewRequest("POST", "/hosts/create", strings.NewReader(requestParams.Encode()))
	if reqErr != nil {
		t.Fatal(reqErr)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := InsertHandler(&models.Host{})
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Wanted response code %v, but got %v", http.StatusOK, rr.Code)
	}

	req, reqErr = http.NewRequest("GET", "/hosts", nil)
	if reqErr != nil {
		t.Fatal(reqErr)
	}

	rr = httptest.NewRecorder()
	handler = IndexHandler(&models.Host{})
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Wanted response code %v, but got %v", http.StatusOK, rr.Code)
	}

	var data interface{}
	decodeError := json.NewDecoder(rr.Body).Decode(&data)
	if decodeError != nil {
		t.Fatal(decodeError)
	}

	if data == nil {
		t.Error("Expected data from hosts endpoint but received nothing")
		t.FailNow()
	}

	parker := data.([]interface{})[0]
	if parker.(map[string]interface{})["first_name"] != "Peter" {
		t.Errorf("Expected firstName to be Peter, not %v", parker.(map[string]interface{})["firstName"])
	}
}
