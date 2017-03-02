package cmd

import "github.com/spf13/cobra"

// RootCmd is the root node in the command tree
var RootCmd = &cobra.Command{
	Use:   "de-vault",
	Short: "Utility for managing Vault for the Discovery Environment",
	Long: `A command-line utility for managing a deployment of Hashicorp's Vault
project. This tool is geared towards CyVerse's Discovery Environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
