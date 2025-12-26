package validator

import (
	"strings"
	"unicode"
)

// ValidateLuhn checks if the provided string passes the Luhn algorithm.
// It ignores spaces and dashes, and returns true if the string is valid.
func ValidateLuhn(input string) bool {
	input = strings.ReplaceAll(input, " ", "")
	if len(input) <= 1 {
		return false
	}
	sum := 0
	for i := len(input) - 1; i >= 0; i-- {
		digitChar := input[i]
		if !unicode.IsDigit(rune(digitChar)) {
			return false // Input must contain only digits
		}
		digit := int(digitChar - '0')
		if (len(input)-1-i)%2 == 1 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}
