package schema

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

type Validator struct {
	Redirect string
}

func New() Validator {
	v := Validator{}
	v.Redirect = redirect

	return v
}

func (v *Validator) ValidateRedirect(data []byte) []error {
	schemaLoader := gojsonschema.NewBytesLoader([]byte(v.Redirect))
	return v.validate(schemaLoader, data)
}

func (v *Validator) validate(loader gojsonschema.JSONLoader, data []byte) []error {
	var err error
	var errors []error

	toCheck := gojsonschema.NewBytesLoader(data)
	result, err := gojsonschema.Validate(loader, toCheck)
	if err != nil {
		return []error{err}
	}
	if !result.Valid() {
		for _, desc := range result.Errors() {
			errors = append(errors, fmt.Errorf("%s", desc.Field()))
		}
		return errors
	}
	return nil
}
