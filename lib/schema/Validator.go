package schema

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
)

type Validator struct {
	schemaLoader gojsonschema.JSONLoader
}

type ValidatorOptions struct {
	Schema string
}

type ValidationResult = gojsonschema.Result
type ValidationError = gojsonschema.ResultError

func NewValidator(opts *ValidatorOptions) (*Validator, error) {
	if opts == nil {
		return nil, fmt.Errorf("NewValidator's options must not be nil")
	}
	if len(opts.Schema) == 0 {
		return nil, fmt.Errorf("Validator's schema must not be empty")
	}
	validator := &Validator{}
	validator.schemaLoader = gojsonschema.NewStringLoader(opts.Schema)
	return validator, nil
}

func (v *Validator) Validate(cfg interface{}) (*ValidationResult, error) {
	if cfg == nil {
		return nil, fmt.Errorf("The configuration object is nil")
	}
	if v.schemaLoader == nil {
		return nil, fmt.Errorf("Validator is not initialized properly")
	}
	documentLoader := gojsonschema.NewGoLoader(cfg)
	return gojsonschema.Validate(v.schemaLoader, documentLoader)
}
