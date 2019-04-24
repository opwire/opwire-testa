package script

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"gopkg.in/yaml.v2"
	"github.com/opwire/opwire-testa/lib/engine"
	"github.com/opwire/opwire-testa/lib/schema"
	"github.com/opwire/opwire-testa/lib/storages"
	"github.com/opwire/opwire-testa/lib/utils"
)

type LoaderOptions interface {}

type Loader struct {
	validator *schema.Validator
	skipInvalidSpecs bool
}

func NewLoader(opts LoaderOptions) (l *Loader, err error) {
	l = new(Loader)
	l.validator, err = schema.NewValidator(&schema.ValidatorOptions{ Schema: scriptSchema })
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (l *Loader) LoadScripts(sourceDirs []string) (map[string]*Descriptor) {
	locators, _ := l.ReadDirs(sourceDirs, ".yml")
	descriptors := l.LoadFiles(locators)
	return descriptors
}

func (l *Loader) LoadFiles(locators []*Locator) (descriptors map[string]*Descriptor) {
	descriptors = make(map[string]*Descriptor, 0)
	for _, locator := range locators {
		descriptors[locator.FullPath] = l.LoadFile(locator)
	}
	return descriptors
}

func (l *Loader) LoadFile(locator *Locator) (*Descriptor) {
	if locator == nil {
		panic(fmt.Errorf("Descriptor must not be nil"))
	}

	// load Test Suite from path
	testsuite := &engine.TestSuite{}

	fs := storages.GetFs()
	file, err1 := fs.Open(locator.FullPath)
	if file != nil {
		defer file.Close()
	}
	if err1 != nil {
		return &Descriptor{
			Locator: locator,
			Error: err1,
		}
	}

	parser := yaml.NewDecoder(file)
	err2 := parser.Decode(testsuite)
	if err2 != nil {
		return &Descriptor{
			Locator: locator,
			Error: err2,
		}
	}

	// validate Test Suite by schema
	result, err3 := l.validator.Validate(testsuite)
	if err3 != nil {
		return &Descriptor{
			Locator: locator,
			TestSuite: testsuite,
			Error: err3,
		}
	}

	if result != nil && !result.Valid() {
		errs := make([]string, len(result.Errors()))
		for i, arg := range result.Errors() {
			errs[i] = arg.String()
		}
		return &Descriptor{
			Locator: locator,
			TestSuite: testsuite,
			Error: utils.CombineErrors("", errs),
		}
	}

	return &Descriptor{
		Locator: locator,
		TestSuite: testsuite,
	}
}

func (l *Loader) ReadDirs(sourceDirs []string, ext string) (locators []*Locator, err error) {
	locators = make([]*Locator, 0)
	for _, sourceDir := range sourceDirs {
		items, err := l.ReadDir(sourceDir, ext)
		if err == nil {
			locators = append(locators, items...)
		}
	}
	return locators, nil
}

func (l *Loader) ReadDir(sourceDir string, ext string) ([]*Locator, error) {
	locators := make([]*Locator, 0)
	err := filepath.Walk(sourceDir, func(path string, f os.FileInfo, err error) error {
		if err == nil && !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				locator := &Locator{}
				locator.FullPath = path
				locator.Home = sourceDir
				locator.Path = strings.TrimPrefix(path, sourceDir)
				locators = append(locators, locator)
			}
		}
		return nil
	})
	return locators, err
}

type Locator struct {
	FullPath string
	Home string
	Path string
}

type Descriptor struct {
	Locator *Locator
	TestSuite *engine.TestSuite
	Error error
}

const scriptSchema string = `{
	"type": "object",
	"properties": {
		"testcases": {
			"type": "array",
			"items": {
				"$ref": "#/definitions/TestCase"
			}
		},
		"skipped": {
			"oneOf": [
				{
					"type": "null"
				},
				{
					"type": "boolean"
				}
			]
		}
	},
	"definitions": {
		"TestCase": {
			"type": "object",
			"properties": {
				"title": {
					"type": "string",
					"minLength": 1
				},
				"version": {
					"oneOf": [
						{
							"type": "null"
						},
						{
							"type": "string"
						}
					]
				},
				"request": {
					"oneOf": [
						{
							"type": "null"
						},
						{
							"$ref": "#/definitions/Request"
						}
					]
				},
				"expectation": {
					"oneOf": [
						{
							"type": "null"
						},
						{
							"$ref": "#/definitions/Expectation"
						}
					]
				},
				"skipped": {
					"oneOf": [
						{
							"type": "null"
						},
						{
							"type": "boolean"
						}
					]
				}
			},
			"additionalProperties": false
		},
		"Request": {
			"type": "object",
			"properties": {
				"method": {
					"type": "string",
					"enum": [ "GET", "PUT", "POST", "PATCH", "DELETE" ]
				},
				"url": {
					"type": "string"
				},
				"pdp": {
					"type": "string"
				},
				"path": {
					"type": "string"
				},
				"headers": {
					"type": "array",
					"items": {
						"type": "object",
						"properties": {
							"name": {
								"type": "string"
							},
							"value": {
								"type": "string"
							}
						}
					}
				},
				"body": {
					"type": "string"
				}
			},
			"additionalProperties": false
		},
		"Expectation": {
			"type": "object",
			"properties": {
				"status-code": {
					"oneOf": [
						{
							"type": "null"
						},
						{
							"type": "object",
							"properties": {
								"is-equal-to": {
									"oneOf": [
										{
											"type": "null"
										},
										{
											"type": "integer"
										}
									]
								},
								"belongs-to": {
									"oneOf": [
										{
											"type": "null"
										},
										{
											"type": "array",
											"items": {
												"type": "integer"
											}
										}
									]
								}
							}
						}
					]
				},
				"headers": {
					"oneOf": [
						{
							"type": "null"
						},
						{
							"type": "object",
							"properties": {
								"has-total": {
									"oneOf": [
										{
											"type": "null"
										},
										{
											"type": "integer"
										}
									]
								},
								"items": {
									"oneOf": [
										{
											"type": "null"
										},
										{
											"type": "array",
											"items": {
												"type": "object",
												"properties": {
													"name": {
														"type": "string"
													},
													"is-equal-to": {
														"oneOf": [
															{
																"type": "null"
															},
															{
																"type": "string"
															}
														]
													}
												}
											}
										}
									]
								}
							}
						}
					]
				},
				"body": {
					"oneOf": [
						{
							"type": "null"
						},
						{
							"type": "object",
							"properties": {
								"has-format": {
									"oneOf": [
										{
											"type": "null"
										},
										{
											"type": "string",
											"enum": ["json", "yaml", "flat"]
										}
									]
								},
								"includes": {
									"oneOf": [
										{
											"type": "null"
										},
										{
											"type": "string"
										}
									]
								},
								"is-equal-to": {
									"oneOf": [
										{
											"type": "null"
										},
										{
											"type": "string"
										}
									]
								},
								"match-with": {
									"oneOf": [
										{
											"type": "null"
										},
										{
											"type": "string"
										}
									]
								}
							}
						}
					]
				}
			}
		}
	},
	"additionalProperties": false
}`
