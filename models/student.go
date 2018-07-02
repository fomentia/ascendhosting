package models

type Student struct {
	FirstName       string
	LastName        string
	CountryOfOrigin string
}

func (s Student) Statement() string {
	return `INSERT INTO students (first_name, last_name, country_of_origin) VALUES ($1, $2, $3)`
}

func (s Student) StatementArgs() []interface{} {
	return []interface{}{s.FirstName, s.LastName, s.CountryOfOrigin}
}

func (s Student) Validations() map[string]Validation {
	return map[string]Validation{
		"FirstName":       lengthGreaterThanZero,
		"LastName":        lengthGreaterThanZero,
		"CountryOfOrigin": lengthGreaterThanZero,
	}
}
