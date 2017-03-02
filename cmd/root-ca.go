package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// root-caCmd represents the root-ca command
var rootCaCmd = &cobra.Command{
	Use:   "root-ca",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("root-ca called")
	},
}

func init() {
	initCmd.AddCommand(rootCaCmd)
	removeCmd.AddCommand(rootCaCmd)
	checkCmd.AddCommand(rootCaCmd)
	rootCaCmd.PersistentFlags().StringVar(
		&mount, // defined in root.go
		"mount",
		"",
		"The path in Vault to the intermediate CA pki backend.",
	)
	rootCaCmd.PersistentFlags().StringVar(
		&role, // defined in root.go
		"role",
		"",
		"The name of the role to use for operations on the intermediate CA.",
	)
	rootCaCmd.PersistentFlags().StringVar(
		&commonName, // defined in root.go
		"common-name",
		"",
		"The common name to use for operations on the intermediate CA.",
	)

}
