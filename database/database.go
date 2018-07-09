package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

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

type Row map[string]interface{}

type Database interface {
	Insert(models.Model) (models.Errors, error)
	Get(models.Model) ([]Row, error)
}

type DB struct {
	sqlDB *sql.DB
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
		stmt := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)",
			model.TableName(),
			model.Columns(),
			valuePlaceholders(model.Columns()))

		_, err := db.sqlDB.Exec(stmt, model.Values()...)
		if err != nil {
			databaseError = err
		}
	}

	return
}

func valuePlaceholders(columnsString string) string {
	columns := strings.Split(columnsString, ", ")
	placeholders := []string{}
	for index, _ := range columns {
		placeholders = append(placeholders, fmt.Sprintf("$%v", index+1))
	}
	return strings.Join(placeholders, ", ")
}

func (db *DB) Get(model models.Model) ([]Row, error) {
	rows := []Row{}

	cursor, err := db.sqlDB.Query(fmt.Sprintf("SELECT %v FROM %v", model.Columns(), model.TableName()))
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
