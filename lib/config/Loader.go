package config

import (
	"github.com/opwire/opwire-testa/lib/schema"
)

type LoaderOptions interface {}

type Loader struct {
	validator *schema.Validator
}

func NewLoader(opts LoaderOptions) (ref *Loader, err error) {
	ref = new(Loader)
	ref.validator, err = schema.NewValidator(&schema.ValidatorOptions{ Schema: configSchema })
	if err != nil {
		return nil, err
	}
	return ref, nil
}

const configSchema string = `{
	"type": "object",
	"properties": {
	}
}`
