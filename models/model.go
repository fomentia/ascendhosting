package models

import (
	"fmt"
	"net/url"
	"strings"
)

type Model interface {
	SetValues(url.Values)
	Get(string) string
	TableName() string
	Columns() string
	Values() []interface{}
	Validations() map[string]Validation
}

type Errors []string

func (e Errors) None() bool {
	return len(e) == 0
}

func (e Errors) Error() string {
	return strings.Join(e, ", ")
}

func Validate(m Model) (errors Errors) {
	for fieldName, validation := range m.Validations() {
		fieldValue := m.Get(fieldName)
		if len(fieldValue) == 0 {
			errors = append(errors, fmt.Sprintf("%v was not supplied", fieldName))
			continue
		}

		if !validation(fieldValue) {
			errors = append(errors, fmt.Sprintf("%v is invalid", fieldName))
		}
	}

	return
}
