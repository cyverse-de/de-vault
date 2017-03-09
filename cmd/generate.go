package cmd

import "github.com/spf13/cobra"

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates the Vault resources represented by the subcommands.",
	Long:  `Generates the Vault resources represented by the subcommands.`,
}

func init() {
	RootCmd.AddCommand(generateCmd)
}
