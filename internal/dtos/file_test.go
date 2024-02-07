package dtos

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

func TestAddNewFileBody_Valid(t *testing.T) {
	testCases := []struct {
		name  string
		body  AddNewFileBody
		valid bool
	}{
		{
			name: "valid",
			body: AddNewFileBody{
				Data: domain.FileData{
					Content: []byte{1, 2, 3},
					Name:    "name",
				},
			},
			valid: true,
		},
		{
			name: "no valid",
			body: AddNewFileBody{
				Data: domain.FileData{
					Content: []byte{},
					Name:    "",
				},
			},
			valid: false,
		},
		{
			name: "no valid content",
			body: AddNewFileBody{
				Data: domain.FileData{
					Content: []byte{},
					Name:    "name",
				},
			},
			valid: false,
		},
		{
			name: "no valid name",
			body: AddNewFileBody{
				Data: domain.FileData{
					Content: []byte{1, 2, 3},
					Name:    "",
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
