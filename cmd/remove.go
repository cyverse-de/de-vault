package cmd

import "github.com/spf13/cobra"

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Clears out the resources associated with the subcommands from Vault.",
	Long:  `Clears out the resources associated with the subcommands from Vault.`,
}

func init() {
	RootCmd.AddCommand(removeCmd)
}
