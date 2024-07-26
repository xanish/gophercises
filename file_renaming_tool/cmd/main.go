package main

import (
	"github.com/xanish/gophercises/file_renaming_tool"
	"path/filepath"
	"regexp"
)

func main() {
	dir, _ := filepath.Abs("./file_renaming_tool/sample")
	file_renaming_tool.Rename(file_renaming_tool.Options{
		Dir:  dir,
		From: *regexp.MustCompile("^(.+?) ([0-9]{4}) [(]([0-9]+) of ([0-9]+)[)][.](.+?)$"),
		To:   "$2 - $1 - $3 of $4.$5",
	})
}
