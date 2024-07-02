package strings_and_bytes

import (
	"strings"
)

func CountWords(word string) int {
	var count int

	if strings.TrimSpace(word) == "" {
		return count
	}

	for i, char := range word {
		str := string(char)
		if str == strings.ToUpper(str) && i > 0 {
			count++
		}
	}

	return count + 1
}
