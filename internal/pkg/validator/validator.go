package validator

// ValidateLuhn checks if the provided string passes the Luhn algorithm.
// It ignores spaces and dashes, and returns true if the string is valid.
func ValidateLuhn(s string) bool {
	var sum int
	var alt bool
	for i := len(s) - 1; i >= 0; i-- {
		char := s[i]
		if char == ' ' || char == '-' {
			continue
		}
		if char < '0' || char > '9' {
			return false
		}
		digit := int(char - '0')
		if alt {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alt = !alt
	}
	return sum%10 == 0 && !alt
}
