package dtos

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

func TestAddNewLogPassBody_Valid(t *testing.T) {
	testCases := []struct {
		name  string
		body  AddNewLogPassBody
		valid bool
	}{
		{
			name: "valid",
			body: AddNewLogPassBody{
				Data: domain.LogPassData{
					Login:    "login",
					Password: "password",
				},
			},
			valid: true,
		},
		{
			name: "no valid",
			body: AddNewLogPassBody{
				Data: domain.LogPassData{
					Login:    "",
					Password: "",
				},
			},
			valid: false,
		},
		{
			name: "no valid login",
			body: AddNewLogPassBody{
				Data: domain.LogPassData{
					Login:    "",
					Password: "password",
				},
			},
			valid: false,
		},
		{
			name: "no valid password",
			body: AddNewLogPassBody{
				Data: domain.LogPassData{
					Login:    "login",
					Password: "",
				},
			},
			valid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			valid := testCase.body.Valid()
			assert.Equal(t, testCase.valid, valid)
		})
	}
}
