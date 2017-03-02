package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var intermediateCAInitCmd = &cobra.Command{
	Use:   "intermediate-ca",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("intermediate-ca init called")
	},
}

var intermediateCACheckCmd = &cobra.Command{
	Use:   "intermediate-ca",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("intermediate-ca check called")
	},
}

var intermediateCARemoveCmd = &cobra.Command{
	Use:   "intermediate-ca",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("intermediate-ca check called")
	},
}

func init() {
	removeCmd.AddCommand(intermediateCARemoveCmd)
	intermediateCARemoveCmd.PersistentFlags().StringVar(
		&mount, // defined in root.go
		"mount",
		"",
		"The path in Vault to the intermediate CA pki backend.",
	)
	intermediateCARemoveCmd.PersistentFlags().StringVar(
		&role, // defined in root.go
		"role",
		"",
		"The name of the role to use for operations on the intermediate CA.",
	)
	intermediateCARemoveCmd.PersistentFlags().StringVar(
		&commonName, // defined in root.go
		"common-name",
		"",
		"The common name to use for operations on the intermediate CA.",
	)

	initCmd.AddCommand(intermediateCAInitCmd)
	intermediateCAInitCmd.PersistentFlags().StringVar(
		&mount, // defined in root.go
		"mount",
		"",
		"The path in Vault to the intermediate CA pki backend.",
	)
	intermediateCAInitCmd.PersistentFlags().StringVar(
		&role, // defined in root.go
		"role",
		"",
		"The name of the role to use for operations on the intermediate CA.",
	)
	intermediateCAInitCmd.PersistentFlags().StringVar(
		&commonName, // defined in root.go
		"common-name",
		"",
		"The common name to use for operations on the intermediate CA.",
	)

	checkCmd.AddCommand(intermediateCACheckCmd)
	intermediateCACheckCmd.PersistentFlags().StringVar(
		&mount, // defined in root.go
		"mount",
		"",
		"The path in Vault to the intermediate CA pki backend.",
	)
	intermediateCACheckCmd.PersistentFlags().StringVar(
		&role, // defined in root.go
		"role",
		"",
		"The name of the role to use for operations on the intermediate CA.",
	)
	intermediateCACheckCmd.PersistentFlags().StringVar(
		&commonName, // defined in root.go
		"common-name",
		"",
		"The common name to use for operations on the intermediate CA.",
	)
}
