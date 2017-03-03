package cmd

import (
	"fmt"
	"log"

	"github.com/cyverse-de/vaulter"
	"github.com/spf13/cobra"
)

var rootCAInitCmd = &cobra.Command{
	Use:   "root-ca",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("root-ca init called")
	},
}

var rootCACheckCmd = &cobra.Command{
	Use:   "root-ca",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Root CA backend is mounted:\t")
		hasRoot, err := vaulter.IsMounted(vaultAPI, mount)
		if err != nil {
			log.Fatal(err)
		}
		if hasRoot {
			fmt.Print("YES\n")
		} else {
			fmt.Print("NO\n")
		}

		var hasRole bool
		fmt.Printf("Root CA role exists:\t")
		hasRole, err = vaulter.HasRole(vaultAPI, mount, role, commonName, true)
		if err != nil {
			log.Fatal(err)
		}
		if hasRole {
			fmt.Print("YES\n")
		} else {
			fmt.Print("NO\n")
		}

		fmt.Printf("Root CA cert exists:\t")
		var hasCert bool
		hasCert, err = vaulter.HasRootCert(vaultAPI, mount, role, commonName)
		if err != nil {
			log.Fatal(err)
		}
		if hasCert {
			fmt.Print("YES\n")
		} else {
			fmt.Print("NO\n")
		}
	},
}

var rootCARemoveCmd = &cobra.Command{
	Use:   "root-ca",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("root-ca remove called")
	},
}

func init() {
	initCmd.AddCommand(rootCAInitCmd)
	rootCAInitCmd.PersistentFlags().StringVar(
		&mount, // defined in root.go
		"mount",
		"",
		"The path in Vault to the intermediate CA pki backend.",
	)
	rootCAInitCmd.PersistentFlags().StringVar(
		&role, // defined in root.go
		"role",
		"",
		"The name of the role to use for operations on the intermediate CA.",
	)
	rootCAInitCmd.PersistentFlags().StringVar(
		&commonName, // defined in root.go
		"common-name",
		"",
		"The common name to use for operations on the intermediate CA.",
	)

	removeCmd.AddCommand(rootCARemoveCmd)
	rootCARemoveCmd.PersistentFlags().StringVar(
		&mount, // defined in root.go
		"mount",
		"",
		"The path in Vault to the intermediate CA pki backend.",
	)
	rootCARemoveCmd.PersistentFlags().StringVar(
		&role, // defined in root.go
		"role",
		"",
		"The name of the role to use for operations on the intermediate CA.",
	)
	rootCARemoveCmd.PersistentFlags().StringVar(
		&commonName, // defined in root.go
		"common-name",
		"",
		"The common name to use for operations on the intermediate CA.",
	)

	checkCmd.AddCommand(rootCACheckCmd)
	rootCACheckCmd.PersistentFlags().StringVar(
		&mount, // defined in root.go
		"mount",
		"",
		"The path in Vault to the intermediate CA pki backend.",
	)
	rootCACheckCmd.PersistentFlags().StringVar(
		&role, // defined in root.go
		"role",
		"",
		"The name of the role to use for operations on the intermediate CA.",
	)
	rootCACheckCmd.PersistentFlags().StringVar(
		&commonName, // defined in root.go
		"common-name",
		"",
		"The common name to use for operations on the intermediate CA.",
	)
}
