package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"reflect"

	_ "github.com/lib/pq"
)

var hostSchema = `CREATE TABLE IF NOT EXISTS hosts (
  id SERIAL PRIMARY KEY,
  first_name VARCHAR(255),
  last_name VARCHAR(255)
)`

type errors []string

func (e *errors) concatenate(delimiter string) string {
	var buffer bytes.Buffer

	for i := 0; i < len(*e); i++ {
		buffer.WriteString((*e)[i])
		if i != len(*e)-1 {
			buffer.WriteString(delimiter)
		}
	}

	return buffer.String()
}

func (e *errors) none() bool {
	return len(*e) == 0
}

func (h *Host) validate() (errors errors) {
	v := reflect.ValueOf(*h)

	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name
		fieldValue := v.Field(i)

		validation, exists := hostValidations[fieldName]
		if exists && validation(fieldValue) != true {
			errors = append(errors, fmt.Sprintf("%v is invalid", fieldName))
		}
	}

	return
}

type validation func(reflect.Value) bool

type Host struct {
	firstName string
	lastName  string
}

var lengthGreaterThanZero = func(data reflect.Value) bool {
	if data.Kind() != reflect.String {
		return false
	}

	return data.Len() != 0
}

var hostValidations = map[string]validation{
	"firstName": lengthGreaterThanZero,
	"lastName":  lengthGreaterThanZero,
}

type DB struct {
	database *sql.DB
}

func initDB() (*DB, error) {
	datastoreName := os.Getenv("ASCEND_HOSTING_POSTGRES_CONNECTION")

	database, err := sql.Open("postgres", datastoreName)
	if err != nil {
		return nil, err
	}

	_, err = database.Exec(hostSchema)
	if err != nil {
		return nil, err
	}

	return &DB{database}, nil
}

func (db *DB) insertHost(host Host) (errors errors) {
	errors = host.validate()

	if errors.none() {
		stmt := `INSERT INTO hosts (first_name, last_name) VALUES ($1, $2)`
		_, err := db.database.Exec(stmt, host.firstName, host.lastName)

		// TODO: separate database errors from validation errors.
		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	return errors
}
