package models

import "net/url"

type Student struct {
	Values url.Values
}

func (s Student) Get(fieldName string) string {
	return s.Values.Get(fieldName)
}

func (s Student) Statement() string {
	return `INSERT INTO students (first_name, last_name, country_of_origin) VALUES ($1, $2, $3)`
}

func (s Student) StatementArgs() []interface{} {
	return []interface{}{s.Values.Get("firstName"), s.Values.Get("lastName"), s.Values.Get("countryOfOrigin")}
}

func (s Student) Validations() map[string]Validation {
	return map[string]Validation{
		"firstName":       lengthGreaterThanZero,
		"lastName":        lengthGreaterThanZero,
		"countryOfOrigin": lengthGreaterThanZero,
	}
}
