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

	rr := post(t, "hosts/create",
		url.Values{"firstName": {"Peter"}, "lastName": {"Parker"}},
		InsertHandler(&models.Host{}))
	if rr.Code != http.StatusOK {
		t.Errorf("Wanted response code %v, but got %v", http.StatusOK, rr.Code)
	}

	var data interface{}
	rr, data = get(t, "hosts", IndexHandler(&models.Host{}))
	if rr.Code != http.StatusOK {
		t.Errorf("Wanted response code %v, but got %v", http.StatusOK, rr.Code)
	}

	parker := data.([]interface{})[0]
	if parker.(map[string]interface{})["first_name"] != "Peter" {
		t.Errorf("Expected firstName to be Peter, not %v", parker.(map[string]interface{})["firstName"])
	}
}

func TestStudents(t *testing.T) {
	DB = &mockDB{}

	rr := post(t, "students/create",
		url.Values{"firstName": {"Tony"}, "lastName": {"Stark"}, "countryOfOrigin": {"U.S.A"}},
		InsertHandler(&models.Student{}))
	if rr.Code != http.StatusOK {
		t.Errorf("Wanted response code %v, but got %v", http.StatusOK, rr.Code)
	}

	var data interface{}
	rr, data = get(t, "students", IndexHandler(&models.Student{}))
	if rr.Code != http.StatusOK {
		t.Errorf("Wanted response code %v, but got %v", http.StatusOK, rr.Code)
	}

	stark := data.([]interface{})[0]
	if stark.(map[string]interface{})["first_name"] != "Tony" {
		t.Errorf("Expected firstName to be Tony, not %v", stark.(map[string]interface{})["firstName"])
	}
}

func post(t *testing.T, url string, body url.Values, handler http.HandlerFunc) httptest.ResponseRecorder {
	req, reqErr := http.NewRequest("POST", url, strings.NewReader(body.Encode()))
	if reqErr != nil {
		t.Fatal(reqErr)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return *rr
}

func get(t *testing.T, url string, handler http.HandlerFunc) (httptest.ResponseRecorder, interface{}) {
	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		t.Fatal(reqErr)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var data interface{}
	decodeError := json.NewDecoder(rr.Body).Decode(&data)
	if decodeError != nil {
		t.Fatal(decodeError)
	}

	if data == nil {
		t.Errorf("Expected data from %v endpoint but received nothing", url)
		t.FailNow()
	}

	return *rr, data
}
