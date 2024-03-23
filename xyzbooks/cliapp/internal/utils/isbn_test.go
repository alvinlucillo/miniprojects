package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertISBN13ToISBN10(t *testing.T) {
	tests := []struct {
		name     string
		isbn13   string
		expected string
		err      error
	}{
		{
			name:     "Valid ISBN-13 with check digit 7",
			isbn13:   "9781234567897",
			expected: "123456789X",
			err:      nil,
		},
		{
			name:     "Valid ISBN 13 with check digit 3",
			isbn13:   "9781891830853",
			expected: "1891830856",
			err:      nil,
		},
		{
			name:     "Valid ISBN 13 with check digit 2",
			isbn13:   "9781603094542",
			expected: "1603094547",
			err:      nil,
		},
		{
			name:     "Invalid ISBN-13",
			isbn13:   "1234567890123",
			expected: "",
			err:      errors.New("invalid ISBN-13"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := ConvertISBN13ToISBN10(test.isbn13)

			if err == nil {
				assert.Nil(t, test.err)
				assert.Equal(t, test.expected, result)
			} else {
				assert.NotNil(t, test.err)
				assert.Equal(t, test.err.Error(), err.Error())
			}
		})
	}
}

func TestConvertISBN10ToISBN13(t *testing.T) {
	tests := []struct {
		name     string
		isbn10   string
		expected string
		err      error
	}{
		{
			name:     "Valid ISBN-10 with check digit X",
			isbn10:   "123456789X",
			expected: "9781234567897",
			err:      nil,
		},
		{
			name:     "Valid ISBN-10 with check digit 6",
			isbn10:   "1891830856",
			expected: "9781891830853",
			err:      nil,
		},
		{
			name:     "Invalid ISBN-10",
			isbn10:   "123456789",
			expected: "",
			err:      errors.New("invalid ISBN-10"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := ConvertISBN10ToISBN13(test.isbn10)

			if err == nil {
				assert.Nil(t, test.err)
				assert.Equal(t, test.expected, result)
			} else {
				assert.NotNil(t, test.err)
				assert.Equal(t, test.err.Error(), err.Error())
			}
		})
	}
}
