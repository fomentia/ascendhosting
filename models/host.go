package models

type Host struct {
	FirstName string
	LastName  string
}

func (h Host) Statement() string {
	return `INSERT INTO hosts (first_name, last_name) VALUES ($1, $2)`
}

func (h Host) StatementArgs() []interface{} {
	return []interface{}{h.FirstName, h.LastName}
}

func (h Host) Validations() map[string]Validation {
	return map[string]Validation{
		"FirstName": lengthGreaterThanZero,
		"LastName":  lengthGreaterThanZero,
	}
}
