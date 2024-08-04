package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xanish/gophercises/secrets_cli_api/vault"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.File(encodingKey, secretsPath())

		key, value := args[0], args[1]
		err := v.Set(key, value)
		if err != nil {
			panic(err)
		}

		fmt.Println("Value set successfully!")
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
}
