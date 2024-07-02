package strings_and_bytes

func Encrypt(message string, rotation int32) string {
	var result []rune

	clamp := func(r, base int32) rune {
		// find out number of positions between rotated char (r) and base (a/A)
		// mod it by 26 to clamp it to only alphabets
		// add the base to get back the character code
		return base + (r-base)%26
	}

	for _, char := range message {
		str := char
		if str >= 'a' && str <= 'z' {
			result = append(result, clamp(char+rotation, 'a'))
		} else if str >= 'A' && str <= 'Z' {
			result = append(result, clamp(char+rotation, 'A'))
		} else {
			result = append(result, char)
		}
	}

	return string(result)
}
