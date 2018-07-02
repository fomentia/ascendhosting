package models

type Validation func(string) bool

var lengthGreaterThanZero = func(data string) bool {
	return len(data) != 0
}
