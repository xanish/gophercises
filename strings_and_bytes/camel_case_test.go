package strings_and_bytes

import "testing"

func TestCountWords(t *testing.T) {
	t.Run("should return 0 if there is just whitespace", func(t *testing.T) {
		want := 0
		got := CountWords("  ")

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("should return 2 for helloWorld", func(t *testing.T) {
		want := 2
		got := CountWords("helloWorld")

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("should return 2 for HelloWorld", func(t *testing.T) {
		want := 2
		got := CountWords("HelloWorld")

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}
