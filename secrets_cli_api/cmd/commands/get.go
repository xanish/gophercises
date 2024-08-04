package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xanish/gophercises/secrets_cli_api/vault"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.File(encodingKey, secretsPath())
		key := args[0]
		value, err := v.Get(key)
		if err != nil {
			fmt.Println("no value set")
			return
		}
		fmt.Printf("%s = %s\n", key, value)
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
