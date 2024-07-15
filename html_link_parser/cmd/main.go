package main

import (
	"fmt"
	"github.com/xanish/gophercises/html_link_parser"
	"os"
	"path/filepath"
)

func main() {
	filePath, _ := filepath.Abs("./html_link_parser/html_link_parser_testdata/ex3.html")
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	parse, err := html_link_parser.Parse(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(parse)
}
