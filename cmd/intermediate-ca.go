package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// intermediate-caCmd represents the intermediate-ca command
var intermediateCaCmd = &cobra.Command{
	Use:   "intermediate-ca",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("intermediate-ca called")
	},
}

func init() {
	initCmd.AddCommand(intermediateCaCmd)
	checkCmd.AddCommand(intermediateCaCmd)
	removeCmd.AddCommand(intermediateCaCmd)

	intermediateCaCmd.PersistentFlags().StringVar(
		&mount, // defined in root.go
		"mount",
		"",
		"The path in Vault to the intermediate CA pki backend.",
	)
	intermediateCaCmd.PersistentFlags().StringVar(
		&role, // defined in root.go
		"role",
		"",
		"The name of the role to use for operations on the intermediate CA.",
	)
	intermediateCaCmd.PersistentFlags().StringVar(
		&commonName, // defined in root.go
		"common-name",
		"",
		"The common name to use for operations on the intermediate CA.",
	)
}
