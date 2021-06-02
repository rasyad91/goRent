package form

import (
	"fmt"
	"net/url"
	"strings"
)

// Form struct
type Form struct {
	url.Values
	Errors errors
}
type errors map[string][]string

// New creates new Form instance
func New(data url.Values) *Form {
	return &Form{data, make(errors)}
}

// Add adds an error message for a given form field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get returns the first error message
func (e errors) Get(field string) string {
	if _, ok := e[field]; !ok || e == nil {
		return ""
	}
	return e[field][0]
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, fmt.Sprintf("%s is mandatory", field))
		}
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) ExistingUser() bool {
	return true
}
