package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/cyverse-de/vaulter"
	vault "github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
)

var rootCAInitCmd = &cobra.Command{
	Use:   "root-ca",
	Short: "Initialize a root CA in Vault",
	Long: `Initializes a root CA in Vault, creating a backend mount, a role, and
a root cert. Requires the --common-name setting. Does not recreate something if
it already exists. If you require a full reset of the mount, role, and/or cert,
use the 'remove root-ca' command followed by a 'init root-ca' command.`,
	Run: func(cmd *cobra.Command, args []string) {
		if mount == "" {
			log.Fatal("--mount must be set.")
		}

		if role == "" {
			log.Fatal("--role must be set.")
		}

		if commonName == "" {
			log.Fatal("--common-name must be set.")
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)

		fmt.Fprint(w, "Mounting root CA backend:\t")
		hasRoot, err := vaulter.IsMounted(vaultAPI, mount)
		if err != nil {
			fmt.Fprintf(w, "FAILURE\t\n")
			log.Fatal(err)
		}
		if !hasRoot {
			if err = vaulter.Mount(vaultAPI, mount, &vaulter.MountConfiguration{
				Type:        "pki",
				MaxLeaseTTL: "87600h",
			}); err != nil {
				fmt.Fprintf(w, "FAILURE\t\n")
				log.Fatal(err)
			}
			fmt.Fprint(w, "SUCCESS\t\n")
		} else {
			fmt.Fprint(w, "SUCCESS\t\n")
		}

		fmt.Fprint(w, "Creating root CA role:\t")
		var hasRole bool
		if hasRole, err = vaulter.HasRole(vaultAPI, mount, role, commonName, true); err != nil {
			fmt.Fprintf(w, "FAILURE\t\n")
			log.Fatal(err)
		}
		if !hasRole {
			_, err = vaulter.CreateRole(vaultAPI, mount, role, &vaulter.RoleConfig{
				AllowedDomains:  commonName,
				AllowSubdomains: true,
				KeyBits:         4096,
				AllowAnyName:    true,
			})
			if err != nil {
				fmt.Fprintf(w, "FAILURE\t\n")
				log.Fatal(err)
			}
			fmt.Fprintf(w, "SUCCESS\t\n")
		} else {
			fmt.Fprintf(w, "SUCCESS\t\n")
		}

		fmt.Fprint(w, "Creating root CA cert:\t")
		var hasCert bool
		if hasCert, err = vaulter.HasRootCert(vaultAPI, mount, role, commonName); err != nil {
			fmt.Fprint(w, "FAILURE\t\n")
			log.Fatal(err)
		}
		if !hasCert {
			var rootCertSecret *vault.Secret
			rootCertSecret, err = vaulter.RootCACert(vaultAPI, mount, &vaulter.RootCACertConfig{
				CommonName: commonName,
				TTL:        "87600h",
				KeyBits:    4096,
			})
			if err != nil {
				fmt.Fprint(w, "FAILURE\t\n")
				log.Fatal(err)
			}
			if rootCertSecret == nil {
				fmt.Fprint(w, "FAILURE\t\n")
				log.Fatal("root CA cert secret is nil")
			}
			fmt.Fprint(w, "SUCCESS\t\n")
		} else {
			fmt.Fprint(w, "SUCCESS\t\n")
		}
		w.Flush()
	},
}

var rootCACheckCmd = &cobra.Command{
	Use:   "root-ca",
	Short: "Checks the status of the root CA in Vault",
	Long: `Checks the status of the root CA in Vault by determining the following:
    1. If the appropriate backend is mounted.
    2. If the role exists.
    3. If the root certificate exists.
This command does not create any of the above if it does not exist. Use the
'init root-ca' command if that is what you require.`,
	Run: func(cmd *cobra.Command, args []string) {
		if mount == "" {
			log.Fatal("--mount must be set.")
		}

		if role == "" {
			log.Fatal("--role must be set.")
		}

		if commonName == "" {
			log.Fatal("--common-name must be set.")
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)

		fmt.Fprintf(w, "Root CA backend is mounted:\t")
		hasRoot, err := vaulter.IsMounted(vaultAPI, mount)
		if err != nil {
			log.Fatal(err)
		}
		if hasRoot {
			fmt.Fprint(w, "YES\t\n")
		} else {
			fmt.Fprint(w, "NO\t\n")
		}

		var hasRole bool
		fmt.Fprintf(w, "Root CA role exists:\t")
		hasRole, err = vaulter.HasRole(vaultAPI, mount, role, commonName, true)
		if err != nil {
			log.Fatal(err)
		}
		if hasRole {
			fmt.Fprint(w, "YES\t\n")
		} else {
			fmt.Fprint(w, "NO\t\n")
		}

		fmt.Fprintf(w, "Root CA cert exists:\t")
		if !hasRole {
			fmt.Fprint(w, "UNKNOWN\t\n")
		} else {
			var hasCert bool
			hasCert, err = vaulter.HasRootCert(vaultAPI, mount, role, commonName)
			if err != nil {
				log.Fatal(err)
			}
			if hasCert {
				fmt.Fprint(w, "YES\t\n")
			} else {
				fmt.Fprint(w, "NO\t\n")
			}
		}

		w.Flush()
	},
}

var rootCARemoveCmd = &cobra.Command{
	Use:   "root-ca",
	Short: "Removes the root CA from Vault",
	Long: `Removes the root CA, the role used with the root CA, and the root cert
from Vault. This is done by unmounting the Vault backend for the root CA. This
command will return successfully if the root CA backend is already unmounted.`,
	Run: func(cmd *cobra.Command, args []string) {
		if mount == "" {
			log.Fatal("--mount must be set.")
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)

		fmt.Fprint(w, "Unmounting root CA backend:\t")
		hasRoot, err := vaulter.IsMounted(vaultAPI, mount)
		if err != nil {
			fmt.Fprintf(w, "FAILURE\t\n")
			log.Fatal(err)
		}
		if hasRoot {
			if err = vaulter.Unmount(vaultAPI, mount); err != nil {
				fmt.Fprintf(w, "FAILURE\t\n")
				log.Fatal(err)
			}
			fmt.Fprintf(w, "SUCCESS\t\n")
		} else {
			fmt.Fprintf(w, "SUCCESS\t\n")
		}
		w.Flush()
	},
}

func init() {
	// Set up the 'init root-ca' command.
	initCmd.AddCommand(rootCAInitCmd)
	rootCAInitCmd.PersistentFlags().StringVar(
		&mount, // defined in root.go
		"mount",
		defaultRootMount,
		"The path in Vault to the intermediate CA pki backend.",
	)
	rootCAInitCmd.PersistentFlags().StringVar(
		&role, // defined in root.go
		"role",
		defaultRootRole,
		"The name of the role to use for operations on the intermediate CA.",
	)
	rootCAInitCmd.PersistentFlags().StringVar(
		&commonName, // defined in root.go
		"common-name",
		"",
		"The common name to use for operations on the intermediate CA.",
	)

	// Set up the 'remove root-ca' command.
	removeCmd.AddCommand(rootCARemoveCmd)
	rootCARemoveCmd.PersistentFlags().StringVar(
		&mount, // defined in root.go
		"mount",
		defaultRootMount,
		"The path in Vault to the intermediate CA pki backend.",
	)

	// Set up the 'check root-ca' command.
	checkCmd.AddCommand(rootCACheckCmd)
	rootCACheckCmd.PersistentFlags().StringVar(
		&mount, // defined in root.go
		"mount",
		defaultRootMount,
		"The path in Vault to the intermediate CA pki backend.",
	)
	rootCACheckCmd.PersistentFlags().StringVar(
		&role, // defined in root.go
		"role",
		defaultRootRole,
		"The name of the role to use for operations on the intermediate CA.",
	)
	rootCACheckCmd.PersistentFlags().StringVar(
		&commonName, // defined in root.go
		"common-name",
		"",
		"The common name to use for operations on the intermediate CA.",
	)
}
