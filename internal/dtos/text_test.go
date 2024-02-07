package dtos

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

func TestAddNewTextBody_Valid(t *testing.T) {
	testCases := []struct {
		name  string
		body  AddNewTextBody
		valid bool
	}{
		{
			name: "valid",
			body: AddNewTextBody{
				Data: domain.TextData{
					Text: "some text",
				},
			},
			valid: true,
		},
		{
			name: "no valid",
			body: AddNewTextBody{
				Data: domain.TextData{
					Text: "",
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
