package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCardNumber(t *testing.T) {
	testCases := []struct {
		name   string
		val    string
		result bool
	}{
		{
			name:   "valid",
			val:    "1234123412341234",
			result: true,
		},
		{
			name:   "valid with spaces",
			val:    "1234 1234 1234 1234",
			result: true,
		},
		{
			name:   "empty",
			val:    "",
			result: false,
		},
		{
			name:   "not valid",
			val:    "123412341234",
			result: false,
		},
		{
			name:   "not valid with spaces",
			val:    "1234 1234 1234  ",
			result: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := ValidateCardNumber(testCase.val)
			assert.Equal(t, testCase.result, result)
		})
	}
}

func TestValidateCVV(t *testing.T) {
	testCases := []struct {
		name   string
		val    string
		result bool
	}{
		{
			name:   "valid",
			val:    "123",
			result: true,
		},
		{
			name:   "empty",
			val:    "",
			result: false,
		},
		{
			name:   "not valid",
			val:    "1234",
			result: false,
		},
		{
			name:   "not valid (letters)",
			val:    "cvv",
			result: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := ValidateCVV(testCase.val)
			assert.Equal(t, testCase.result, result)
		})
	}
}

func TestValidateExpiredAt(t *testing.T) {
	testCases := []struct {
		name   string
		val    string
		result bool
	}{
		{
			name:   "valid",
			val:    "12/30",
			result: true,
		},
		{
			name:   "valid",
			val:    "03/30",
			result: true,
		},
		{
			name:   "not valid (invalid separator)",
			val:    "12.30",
			result: false,
		},
		{
			name:   "empty",
			val:    "",
			result: false,
		},
		{
			name:   "not valid (invalid month)",
			val:    "14/30",
			result: false,
		},
		{
			name:   "not valid (letters)",
			val:    "ff/ff",
			result: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := ValidateExpiredAt(testCase.val)
			assert.Equal(t, testCase.result, result)
		})
	}
}
