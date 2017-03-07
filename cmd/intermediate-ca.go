package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/cyverse-de/vaulter"
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
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)

		fmt.Fprint(w, "Intermediate CA backend is mounted:\t")
		hasIntermediate, err := vaulter.IsMounted(vaultAPI, mount)
		if err != nil {
			log.Fatal(err)
		}
		if hasIntermediate {
			fmt.Fprint(w, "YES\t\n")
		} else {
			fmt.Fprint(w, "NO\t\n")
		}

		fmt.Fprint(w, "Intermediate CA role exists:\t")
		if !hasIntermediate {
			fmt.Fprint(w, "UNKNOWN\t\n")
		} else {
			hasRole, err := vaulter.HasRole(vaultAPI, mount, role, commonName, true)
			if err != nil {
				log.Fatal(err)
			}
			if hasRole {
				fmt.Fprint(w, "YES\t\n")
			} else {
				fmt.Fprint(w, "NO\t\n")
			}
		}

		fmt.Fprint(w, "Intermediate CA backend is configured correctly:\t")
		if !hasIntermediate {
			fmt.Fprint(w, "UNKNOWN\t\n")
		} else {
			config, err := vaulter.ReadMount(vaultAPI, fmt.Sprintf("%s/config/urls", mount), parentToken)
			if err != nil {
				log.Fatal(err)
			}
			if config == nil {
				log.Fatal("secret was nil")
			}
			var (
				v  string
				ok bool
			)
			if v, ok = config["issuing_certificates"].(string); !ok {
				log.Fatal("issuing_certificates was not found")
			}
			if v != fmt.Sprintf("%s/v1/%s/ca", vaultURL, mount) {
				fmt.Fprint(w, "FAILURE\t\n")
				log.Fatalf("issuing_certificates was %s", v)
			}

			if v, ok = config["crl_distribution_points"].(string); !ok {
				fmt.Fprint(w, "FAILURE\t\n")
				log.Fatal("crl_distribution_points was not found")
			}
			if v != fmt.Sprintf("%s/v1/%s/crl", vaultURL, mount) {
				fmt.Fprintf(w, "FAILURE\t\n")
				log.Fatalf("crl_distribution_points was %s", v)
			}
		}
		w.Flush()
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
