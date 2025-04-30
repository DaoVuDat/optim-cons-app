package util

import (
	"fmt"
	"regexp"
)

// IsTFNumber checks if a string follows the pattern "TF<number>" or "tf<number>"
func IsTFNumber(s string) bool {
	re := regexp.MustCompile(`^(?i)(tf)\d+$`)
	return re.MatchString(s)
}

// ValidateTFNumber checks if a string follows the pattern "TF<number>" or "tf<number>" and returns an error if it doesn't
func ValidateTFNumber(s string) error {
	if !IsTFNumber(s) {
		return fmt.Errorf("invalid format: %s must be in the format TF<number> or tf<number>", s)
	}
	return nil
}
