package models

import "net/url"

type Student struct {
	data url.Values
}

func (s *Student) Init(values url.Values) {
	s.data = values
}

func (s *Student) Get(fieldName string) string {
	return s.data.Get(fieldName)
}

func (s *Student) TableName() string {
	return "students"
}

func (s *Student) Columns() string {
	return `first_name, last_name, country_of_origin`
}

func (s *Student) Values() []interface{} {
	return []interface{}{s.data.Get("firstName"), s.data.Get("lastName"), s.data.Get("countryOfOrigin")}
}

func (s *Student) Validations() map[string]Validation {
	return map[string]Validation{
		"firstName":       lengthGreaterThanZero,
		"lastName":        lengthGreaterThanZero,
		"countryOfOrigin": lengthGreaterThanZero,
	}
}
