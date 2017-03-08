package cmd

import "github.com/spf13/cobra"

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the resources represented by the subcommands.",
	Long:  `Initializes the resources represented by the subcommands.`,
}

func init() {
	RootCmd.AddCommand(initCmd)
}
