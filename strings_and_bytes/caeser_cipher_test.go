package strings_and_bytes

import "testing"

func TestEncrypt(t *testing.T) {
	t.Run("should return 'Wkhuh'v-d-vwdupdq-zdlwlqj-lq-wkh-vnb' for 'There's-a-starman-waiting-in-the-sky'", func(t *testing.T) {
		want := "Wkhuh'v-d-vwdupdq-zdlwlqj-lq-wkh-vnb"
		got := Encrypt("There's-a-starman-waiting-in-the-sky", 3)

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
