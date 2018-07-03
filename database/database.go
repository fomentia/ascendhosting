package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/fomentia/ascendhosting/models"
)

// TODO: Separate schema from operations.
var schema = `CREATE TABLE IF NOT EXISTS hosts (
  id SERIAL PRIMARY KEY,
  first_name VARCHAR(255) NOT NULL,
  last_name  VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS students (
  id SERIAL PRIMARY KEY,
  first_name        VARCHAR(255) NOT NULL,
  last_name         VARCHAR(255) NOT NULL,
  country_of_origin VARCHAR(255) NOT NULL
);`

type DB struct {
	Database *sql.DB
}

func InitDB() (*DB, error) {
	datastoreName := os.Getenv("ASCEND_HOSTING_POSTGRES_CONNECTION")

	database, err := sql.Open("postgres", datastoreName)
	if err != nil {
		return nil, err
	}

	_, err = database.Exec(schema)
	if err != nil {
		return nil, err
	}

	return &DB{database}, nil
}

func (db *DB) Insert(model models.Model) (validationErrors models.Errors, databaseError error) {
	validationErrors = models.Validate(model)

	if validationErrors.None() {
		_, err := db.Database.Exec(model.Statement(), model.StatementArgs()...)
		if err != nil {
			databaseError = err
		}
	}

	return
}

type Row map[string]interface{}

func (db *DB) Get(columns string, tableName string) ([]Row, error) {
	rows := []Row{}

	cursor, err := db.Database.Query(fmt.Sprintf("SELECT %v FROM %v", columns, tableName))
	if err != nil {
		return rows, err
	}

	for cursor.Next() {
		row := Row{}

		columns, err := cursor.Columns()
		if err != nil {
			return rows, err
		}

		values := make([]string, len(columns))
		pointers := make([]interface{}, len(columns))

		for index, _ := range values {
			pointers[index] = &values[index]
		}

		scanErr := cursor.Scan(pointers...)
		if scanErr != nil {
			return rows, err
		}

		for index, column := range columns {
			row[column] = values[index]
		}

		rows = append(rows, row)
	}
	if cursor.Err() != nil {
		return rows, err
	}

	return rows, nil
}
