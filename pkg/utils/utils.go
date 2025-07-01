package utils

func RemoveSymbols(s string) string {
	// Create a new string builder to build the result
	var result []rune

	// Iterate over each character in the input string
	for _, char := range s {
		// Check if the character is a digit or a letter
		if (char >= '0' && char <= '9') || (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			result = append(result, char)
		}
	}

	// Convert the rune slice back to a string and return it
	return string(result)
}

func IsNullOrWhiteSpace(s string) bool {
	if s == "" {
		return true
	}
	for _, r := range s {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			return false
		}
	}
	return true
}
