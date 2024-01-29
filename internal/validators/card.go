package validators

import (
	"strconv"
	"strings"
)

// ValidateCardNumber - remove spaces in given string and check if length of string is equal 16
func ValidateCardNumber(number string) bool {
	clearNumber := strings.ReplaceAll(number, " ", "")

	if len(clearNumber) != 16 {
		return false
	}

	return true
}

// ValidateExpiredAt - validate date in format mm/yy (e.g 04/30)
func ValidateExpiredAt(expiredAt string) bool {
	parts := strings.Split(expiredAt, "/")
	if len(parts) != 2 {
		return false
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil {
		return false
	}

	_, err = strconv.Atoi(parts[1])
	if err != nil {
		return false
	}

	if month <= 0 || month > 12 {
		return false
	}

	return true
}

// ValidateCVV - check if cvv len is 3 and it is a valid number
func ValidateCVV(cvv string) bool {
	if len(cvv) != 3 {
		return false
	}

	_, err := strconv.Atoi(cvv)
	if err != nil {
		return false
	}

	return true
}
