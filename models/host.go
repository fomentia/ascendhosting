package models

import "net/url"

type Host struct {
	Values url.Values
}

func (h Host) Get(fieldName string) string {
	return h.Values.Get(fieldName)
}

func (h Host) Statement() string {
	return `INSERT INTO hosts (first_name, last_name) VALUES ($1, $2)`
}

func (h Host) StatementArgs() []interface{} {
	return []interface{}{h.Values.Get("firstName"), h.Values.Get("lastName")}
}

func (h Host) Validations() map[string]Validation {
	return map[string]Validation{
		"firstName": lengthGreaterThanZero,
		"lastName":  lengthGreaterThanZero,
	}
}
