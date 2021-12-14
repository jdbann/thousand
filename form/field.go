package form

import "regexp"

type stringField struct {
	Message string
	Value   string
}

type stringValidations []stringValidation

func (vs stringValidations) validate(f *stringField) bool {
	for _, v := range vs {
		if !v(f) {
			return false
		}
	}
	return true
}

type stringValidation func(field *stringField) bool

func stringPresent(msg string) stringValidation {
	return func(f *stringField) bool {
		if f.Value == "" {
			f.Message = msg
			return false
		}

		return true
	}
}

var emailRegex = regexp.MustCompile(`(?:[a-z0-9!#$%&'*+/=?^_\x60{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_\x60{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])`)

func stringEmailFormat(msg string) stringValidation {
	return func(f *stringField) bool {
		if !emailRegex.Match([]byte(f.Value)) {
			f.Message = msg
			return false
		}

		return true
	}
}
