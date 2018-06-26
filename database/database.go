package database

import (
	"database/sql"
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
	database *sql.DB
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
		err := model.Insert(db.database)
		if err != nil {
			databaseError = err
		}
	}

	return
}
