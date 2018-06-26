package models

import "database/sql"

type Student struct {
	FirstName       string
	LastName        string
	CountryOfOrigin string
}

func (s Student) Insert(db *sql.DB) error {
	stmt := `INSERT INTO students (first_name, last_name, country_of_origin) VALUES ($1, $2, $3)`
	_, err := db.Exec(stmt, s.FirstName, s.LastName, s.CountryOfOrigin)
	return err
}

func (s Student) Validations() map[string]Validation {
	return map[string]Validation{
		"FirstName":       lengthGreaterThanZero,
		"LastName":        lengthGreaterThanZero,
		"CountryOfOrigin": lengthGreaterThanZero,
	}
}
