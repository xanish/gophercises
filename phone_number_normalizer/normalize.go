package phone_number_normalizer

import "regexp"

func Normalize(phone string) string {
	// set regex to track all non-digit characters
	re := regexp.MustCompile("\\D")

	return re.ReplaceAllString(phone, "")
}
