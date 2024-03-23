package utils

import (
	"errors"
	"strconv"
)

// Converts ISBN-13 to ISBN-10
func ConvertISBN13ToISBN10(isbn13 string) (string, error) {
	if len(isbn13) != 13 || (isbn13[:3] != "978" && isbn13[:3] != "979") {
		return "", errors.New("invalid ISBN-13")
	}

	isbn10 := isbn13[3:12]
	sum := 0
	for i, digit := range isbn10 {
		num, err := strconv.Atoi(string(digit))
		if err != nil {
			return "", err
		}
		sum += num * (10 - i)
	}

	checkDigit := 11 - (sum % 11)
	if checkDigit == 10 {
		isbn10 += "X"
	} else if checkDigit == 11 {
		isbn10 += "0"
	} else {
		isbn10 += strconv.Itoa(checkDigit)
	}

	return isbn10, nil
}

// Converts ISBN-10 to ISBN-13
func ConvertISBN10ToISBN13(isbn10 string) (string, error) {
	if len(isbn10) != 10 {
		return "", errors.New("invalid ISBN-10")
	}

	isbn13 := "978" + isbn10[:9]
	sum := 0
	for i, digit := range isbn13 {
		num, err := strconv.Atoi(string(digit))
		if err != nil {
			return "", err
		}
		if i%2 == 0 {
			sum += num
		} else {
			sum += num * 3
		}
	}

	checkDigit := 10 - (sum % 10)
	if checkDigit == 10 {
		checkDigit = 0
	}
	isbn13 += strconv.Itoa(checkDigit)

	return isbn13, nil
}
