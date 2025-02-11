package tools

import "unicode"

func IsCharInRange(r rune) bool {
	if unicode.Is(unicode.Latin, r) {
		return true
	}

	if unicode.Is(unicode.Cyrillic, r) {
		return true
	}

	if unicode.IsDigit(r) {
		return true
	}

	switch r {
	case ',', '*', '!', ')', '(', '#', '@', '$', '%', '&', '<', '>', '"', '\'', '\\', '|', '/', ':', ';', '^', '?', '-', '_', '=', '+', '.':
		return true
	}

	return false
}
