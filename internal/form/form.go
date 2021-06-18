package form

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
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

// CheckLength validates the length of the input string, if there is no max length, enter negative integer (eg. -1) for max field
func (f *Form) CheckLength(field string, min, max int) {
	value := f.Get(field)
	if max < 0 {
		if len(value) < min {
			f.Errors.Add(field, fmt.Sprintf("%s min. length is %d", field, min))
		}
	} else {
		if len(value) < min || len(value) > max {
			f.Errors.Add(field, fmt.Sprintf("%s min.length: %d and max.length: %d", field, min, max))
		}
	}
}

//CheckEmail validates the email
func (f *Form) CheckEmail(field string) {
	value := f.Get(field)
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	validEmail := emailRegex.MatchString(value)
	if !validEmail {
		f.Errors.Add(field, fmt.Sprintf("Please enter valid email"))
	}
}

// Required checks whether the fields are populated
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, fmt.Sprintf("%s is mandatory", field))
		}
	}
}

// Valid will return false if there is any errors in Form
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

//retrieves category informatin
func (f *Form) RetrieveCategory(field string) string {
	var startWordIndex int
	value := strings.TrimSpace(f.Get(field))

	categoryArray := strings.Split(value, "")

	lastIndex := len(categoryArray) - 1

	for p := lastIndex; p > -1; p-- {
		if categoryArray[p] == ">" {
			startWordIndex = p + 2
			break
		}
	}
	categoryArray = categoryArray[startWordIndex:]
	res := strings.Join(categoryArray, "")
	return res
}

//retrieves category informatin
func (f *Form) ProcessPrice(field string) float32 {

	value := f.Get(field)
	res, err := strconv.ParseFloat(value, 64)
	if err != nil {
		f.Errors.Add(field, fmt.Sprintf("%s should not contain alphabets or any letters", field))
	}

	return float32(res)
}
