package html_link_parser

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	files := []struct {
		file string
		want []Link
	}{
		{
			file: "./html_link_parser_testdata/ex1.html",
			want: []Link{{"/other-page", "A link to another page"}},
		},
		{
			file: "./html_link_parser_testdata/ex2.html",
			want: []Link{{
				"https://www.twitter.com/joncalhoun",
				"Check me out on twitter",
			}, {
				"https://github.com/gophercises",
				"Gophercises is onGithub!",
			}},
		},
		{
			file: "./html_link_parser_testdata/ex3.html",
			want: []Link{{
				"#",
				"Login",
			}, {
				"/lost",
				"Lost? Need help?",
			}, {
				"https://twitter.com/marcusolsson",
				"@marcusolsson",
			}},
		},
		{
			file: "./html_link_parser_testdata/ex4.html",
			want: []Link{{"/dog-cat", "dog cat"}},
		},
	}

	for _, tt := range files {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run("parse from "+tt.file, func(t *testing.T) {
			filePath, _ := filepath.Abs(tt.file)
			f, err := os.Open(filePath)
			if err != nil {
				t.Fatalf("could not open file %s", filePath)
			}

			parse, err := Parse(f)
			if err != nil {
				t.Fatalf("could not parse file %s", filePath)
			}

			if reflect.DeepEqual(parse, tt.want) == false {
				t.Errorf("expected: %v, got: %v", tt.want, parse)
			}
		})
	}
}
