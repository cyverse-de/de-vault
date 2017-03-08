package cmd

import "github.com/spf13/cobra"

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Clears out the resources associated with the subcommands from Vault.",
	Long:  `Clears out the resources associated with the subcommands from Vault.`,
	// Run: func(cmd *cobra.Command, args []string) {
	//
	// },
}

func init() {
	RootCmd.AddCommand(removeCmd)
}
