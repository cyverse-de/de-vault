package cmd

import "github.com/spf13/cobra"

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Reports the status of the Vault resources represented by the subcommands.",
	Long:  `Reports the status of the Vault resources represented by the subcommands.`,
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
