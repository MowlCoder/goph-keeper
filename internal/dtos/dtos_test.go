package dtos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteBatchBody_Valid(t *testing.T) {
	testCases := []struct {
		name  string
		body  DeleteBatchBody
		valid bool
	}{
		{
			name: "valid",
			body: DeleteBatchBody{
				IDs: []int{1, 2, 3},
			},
			valid: true,
		},
		{
			name: "no valid (empty)",
			body: DeleteBatchBody{
				IDs: []int{},
			},
			valid: false,
		},
		{
			name: "no valid (nil)",
			body: DeleteBatchBody{
				IDs: nil,
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
