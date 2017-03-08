package cmd

import "github.com/spf13/cobra"

var revokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revokes the Vault resource represented by the subcommand.",
	Long:  `Revokes the Vault resource represented by the subcommand.`,
}

func init() {
	RootCmd.AddCommand(revokeCmd)
}
