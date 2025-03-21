package util

import (
	"errors"
	"regexp"
)

var invalidInputError = errors.New("invalid input")

func IsValidPositiveInteger(input string) error {
	re := regexp.MustCompile(`^[1-9]\d*$`)

	if re.MatchString(input) {
		return nil
	} else {
		return invalidInputError
	}
}

func IsValidFloat(input string) error {
	re := regexp.MustCompile(`^[1-9]\d*(\.\d+)?$`)

	if re.MatchString(input) {
		return nil
	} else {
		return invalidInputError
	}
}

func IsValidBoolean(input string) error {
	re := regexp.MustCompile(`^(true|false)$`)

	if re.MatchString(input) {
		return nil
	} else {
		return invalidInputError
	}
}

func IsValidFList(input string) error {
	re := regexp.MustCompile(`^[0-9]+(\.[0-9]+)?(,[0-9]+(\.[0-9]+)?)*$`)

	if re.MatchString(input) {
		return nil
	} else {
		return invalidInputError
	}
}
