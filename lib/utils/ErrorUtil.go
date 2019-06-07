package utils

import (
	"errors"
	"fmt"
	"strings"
)

func BuildMultilineError(errs []string) error {
	if len(errs) == 0 {
		return nil
	}
	return errors.New(strings.Join(errs, "\n"))
}

func CombineErrors(label string, messages []string) error {
	errstrs := make([]string, 0)
	for _, msg := range messages {
		if len(msg) > 0 {
			errstrs = append(errstrs, msg)
		}
	}
	if len(errstrs) > 0 {
		errstrs = append([]string {label}, errstrs...)
		return fmt.Errorf(strings.Join(errstrs, "\n - "))
	}
	return nil
}

func LabelifyError(label string, err error) error {
	lines := strings.Split(err.Error(), "\n")
	lines = AppendLinesWithIndent([]string{label}, lines, 2)
	return BuildMultilineError(lines)
}
