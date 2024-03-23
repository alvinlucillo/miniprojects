package utils

import "testing"

func TestValidateISBN13(t *testing.T) {
	tests := []struct {
		isbn13 string
		want   bool
	}{
		{"9783161484100", true},
		{"978316148410", false},
		{"97831614841001", false},
		{"978316148410a", false},
		{"978316148410A", false},
		{"978316148410!", false},
	}

	for _, tt := range tests {
		got := ValidateISBN13(tt.isbn13)
		if got != tt.want {
			t.Errorf("ValidateISBN13(%s) = %v, want %v", tt.isbn13, got, tt.want)
		}
	}
}

func TestValidateISBN10(t *testing.T) {
	tests := []struct {
		isbn10 string
		want   bool
	}{
		{"0471958697", true},
		{"0471958698", false},
		{"1234567890", false},
		{"abcabcabca", false},
	}

	for _, tt := range tests {
		got := ValidateISBN10(tt.isbn10)
		if got != tt.want {
			t.Errorf("ValidateISBN10(%s) = %v, want %v", tt.isbn10, got, tt.want)
		}
	}
}
