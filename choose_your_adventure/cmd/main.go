package main

import (
	"errors"
	"fmt"
	"github.com/xanish/gophercises/choose_your_adventure"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const serveAddr = ":8080"

func main() {
	filePath, _ := filepath.Abs("./choose_your_adventure/choose_your_adventure_testdata/gopher.json")
	storyJson, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open story file: %v", err)
	}

	defer func(storyJson *os.File) {
		_ = storyJson.Close()
	}(storyJson)

	filePath, err = filepath.Abs("./choose_your_adventure/templates/main.gohtml")
	if err != nil {
		log.Fatalf("failed to open template file: %v", err)
	}

	//createAndExecuteHTTP(storyJson, filePath)
	createAndExecuteCLI(storyJson, filePath)
}

func createAndExecuteHTTP(storyJson io.Reader, filePath string) {
	handler, err := choose_your_adventure.Web(storyJson, filePath)
	if err != nil {
		log.Fatalf("failed to create web handler: %v", err)
	}

	err = http.ListenAndServe(serveAddr, handler)
	if err != nil {
		log.Fatalf("could not start http server: %v", err)
	}

	fmt.Printf("server started on port %s", serveAddr)
}

func createAndExecuteCLI(storyJson io.Reader, filePath string) {
	cli, err := choose_your_adventure.CLI(storyJson, os.Stdin)
	if err != nil {
		return
	}

	err = cli.ServeCLI("", os.Stdout)
	if errors.Is(err, choose_your_adventure.ErrAdventureEnded) {
		return
	} else {
		log.Fatalf("failed to start cli app: %v", err)
	}
}
