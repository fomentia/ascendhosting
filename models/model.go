package models

import (
	"bytes"
	"fmt"
)

type Model interface {
	Get(string) string
	Statement() string
	StatementArgs() []interface{}
	Validations() map[string]Validation
}

type Errors []string

func (e *Errors) Concatenate(delimiter string) string {
	var buffer bytes.Buffer

	for i := 0; i < len(*e); i++ {
		buffer.WriteString((*e)[i])
		if i != len(*e)-1 {
			buffer.WriteString(delimiter)
		}
	}

	return buffer.String()
}

func (e *Errors) None() bool {
	return len(*e) == 0
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
