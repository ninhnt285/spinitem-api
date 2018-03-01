package validation

import (
	"regexp"
	"strings"
)

// ValidateEmail validates email input
func ValidateEmail(email string) bool {
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, email); !m {
		return false
	}
	return true
}

// ValidateStringInput validates string input is not blank
func ValidateStringInput(username string) bool {
	if strings.TrimSpace(username) == "" {
		return false
	}
	return true
}
