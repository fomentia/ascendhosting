package models

import "net/url"

type Host struct {
	data url.Values
}

func (h *Host) SetValues(values url.Values) {
	h.data = values
}

func (h *Host) Get(fieldName string) string {
	return h.data.Get(fieldName)
}

func (h *Host) TableName() string {
	return "hosts"
}

func (h *Host) Columns() string {
	return `first_name, last_name`
}

func (h *Host) Values() []interface{} {
	return []interface{}{h.Get("firstName"), h.Get("lastName")}
}

func (h *Host) Validations() map[string]Validation {
	return map[string]Validation{
		"firstName": lengthGreaterThanZero,
		"lastName":  lengthGreaterThanZero,
	}
}
