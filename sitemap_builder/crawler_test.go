package sitemap_builder

import (
	"bytes"
	"reflect"
	"testing"
)

func TestCrawl(t *testing.T) {
	buf := bytes.Buffer{}
	want := []string{"https://gophercises.com", "https://gophercises.com/demos/cyoa"}
	got, err := Crawl("https://gophercises.com", 2, &buf)
	if err != nil {
		t.Fatalf("Error building sitemap %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v got %v", want, got)
	}
}
