package utils

func ValidateISBN13(isbn13 string) bool {
	if len(isbn13) != 13 {
		return false
	}

	sum := 0
	for i := 0; i < 12; i++ {
		digit := int(isbn13[i] - '0')
		if i%2 == 0 {
			sum += digit
		} else {
			sum += 3 * digit
		}
	}

	check := (10 - (sum % 10)) % 10
	return check == int(isbn13[12]-'0')
}

func ValidateISBN10(isbn10 string) bool {
	if len(isbn10) != 10 {
		return false
	}

	sum := 0
	for i := 0; i < 9; i++ {
		digit := int(isbn10[i] - '0')
		sum += (i + 1) * digit
	}

	check := sum % 11
	if check == 10 {
		return isbn10[9] == 'X'
	}

	return check == int(isbn10[9]-'0')
}
