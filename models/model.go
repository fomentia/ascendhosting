package models

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
)

type Model interface {
	Insert(*sql.DB) error
	Validations() map[string]Validation
}

type Validation func(reflect.Value) bool

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
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Struct {
		errors = append(errors, "Model is not a struct")
		return
	}

	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name
		fieldValue := v.Field(i)

		validation, exists := m.Validations()[fieldName]
		if exists && validation(fieldValue) != true {
			errors = append(errors, fmt.Sprintf("%v is invalid", fieldName))
		}
	}

	return
}
