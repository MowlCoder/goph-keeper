package dtos

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

func TestAddNewCardBody_Valid(t *testing.T) {
	testCases := []struct {
		name  string
		body  AddNewCardBody
		valid bool
	}{
		{
			name: "valid",
			body: AddNewCardBody{
				Data: domain.CardData{
					Number:    "1234123412341234",
					ExpiredAt: "12/30",
					CVV:       "123",
				},
			},
			valid: true,
		},
		{
			name: "no valid",
			body: AddNewCardBody{
				Data: domain.CardData{
					Number:    "",
					ExpiredAt: "",
					CVV:       "",
				},
			},
			valid: false,
		},
		{
			name: "no valid number",
			body: AddNewCardBody{
				Data: domain.CardData{
					Number:    "",
					ExpiredAt: "12/30",
					CVV:       "123",
				},
			},
			valid: false,
		},
		{
			name: "no valid expired at",
			body: AddNewCardBody{
				Data: domain.CardData{
					Number:    "1234123412341234",
					ExpiredAt: "",
					CVV:       "123",
				},
			},
			valid: false,
		},
		{
			name: "no valid cvv",
			body: AddNewCardBody{
				Data: domain.CardData{
					Number:    "1234123412341234",
					ExpiredAt: "12/30",
					CVV:       "",
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
