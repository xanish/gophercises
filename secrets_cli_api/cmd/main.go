package main

import (
	"fmt"

	"github.com/xanish/gophercises/secrets_cli_api/vault"
)

func main() {
	v := vault.Memory("my-fake-key")

	err := v.Set("demo_key", "some crazy value")
	if err != nil {
		panic(err)
	}

	plain, err := v.Get("demo_key")
	if err != nil {
		panic(err)
	}

	fmt.Println("Plain:", plain)
}
