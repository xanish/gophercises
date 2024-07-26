package file_renaming_tool

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

type Options struct {
	Dir    string
	From   regexp.Regexp
	To     string
	DryRun bool
}

type renameResult struct {
	from string
	to   string
}

func Rename(opts Options) error {
	var toRename []renameResult
	err := filepath.Walk(opts.Dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if opts.From.MatchString(info.Name()) {
			newName := opts.From.ReplaceAllString(info.Name(), opts.To)
			toRename = append(toRename, renameResult{path, filepath.Join(filepath.Dir(path), newName)})
		}

		return nil
	})

	if err != nil {
		return err
	}

	for _, rename := range toRename {
		fmt.Printf("mv %s => %s\n", rename.from, rename.to)
		if !opts.DryRun {
			err = os.Rename(rename.from, rename.to)
			if err != nil {
				fmt.Printf("error renaming %s: %v", rename.from, err)
			}
		}
	}

	return nil
}
