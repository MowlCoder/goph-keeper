package dtos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterBody_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		body  RegisterBody
		valid bool
	}{
		{
			name: "valid",
			body: RegisterBody{
				Email:    "email@email.com",
				Password: "Test1+",
			},
			valid: true,
		},
		{
			name: "no valid",
			body: RegisterBody{
				Email:    "",
				Password: "",
			},
			valid: false,
		},
		{
			name: "no valid email",
			body: RegisterBody{
				Email:    "email",
				Password: "Test1+",
			},
			valid: false,
		},
		{
			name: "no valid password",
			body: RegisterBody{
				Email:    "email@email.com",
				Password: "Test11111",
			},
			valid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			valid := testCase.body.Validate()
			assert.Equal(t, testCase.valid, valid)
		})
	}
}

func TestAuthorizeBody_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		body  AuthorizeBody
		valid bool
	}{
		{
			name: "valid",
			body: AuthorizeBody{
				Email:    "email@email.com",
				Password: "Password",
			},
			valid: true,
		},
		{
			name: "no valid",
			body: AuthorizeBody{
				Email:    "",
				Password: "",
			},
			valid: false,
		},
		{
			name: "no valid email",
			body: AuthorizeBody{
				Email:    "",
				Password: "Password",
			},
			valid: false,
		},
		{
			name: "no valid password",
			body: AuthorizeBody{
				Email:    "email@email.com",
				Password: "",
			},
			valid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			valid := testCase.body.Validate()
			assert.Equal(t, testCase.valid, valid)
		})
	}
}
