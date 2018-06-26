package models

import "database/sql"

type Host struct {
	FirstName string
	LastName  string
}

func (h Host) Insert(db *sql.DB) error {
	stmt := `INSERT INTO hosts (first_name, last_name) VALUES ($1, $2)`
	_, err := db.Exec(stmt, h.FirstName, h.LastName)
	return err
}

func (h Host) Validations() map[string]Validation {
	return map[string]Validation{
		"FirstName": lengthGreaterThanZero,
		"LastName":  lengthGreaterThanZero,
	}
}
