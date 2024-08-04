package main

import "github.com/xanish/gophercises/secrets_cli_api/cmd/commands"

func main() {
	err := commands.RootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
