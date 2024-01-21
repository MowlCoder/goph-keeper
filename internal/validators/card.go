package validators

import (
	"strconv"
	"strings"
)

func ValidateCardNumber(number string) bool {
	clearNumber := strings.ReplaceAll(number, " ", "")

	if len(clearNumber) != 16 {
		return false
	}

	return true
}

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
