package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"text/tabwriter"

	"github.com/cyverse-de/vaulter"
	"github.com/spf13/cobra"
)

var (
	rootmount string
)

var intermediateCAInitCmd = &cobra.Command{
	Use:   "intermediate-ca",
	Short: "Initialize an intermediate CA in Vault.",
	Long: `Initializes an intermediate CA in Vault, the end result being a new PKI
 backend that has a role configured and a signed CSR imported into it.`,
	Run: func(cmd *cobra.Command, args []string) {
		if mount == "" {
			log.Fatal("--mount was not set.")
		}
		if role == "" {
			log.Fatal("--role was not set.")
		}
		if commonName == "" {
			log.Fatal("--common-name was not set.")
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)

		fmt.Fprintf(w, "Creating the intermediate CA:\t")
		hasIntermediate, err := vaulter.IsMounted(vaultAPI, mount)
		if err != nil {
			fmt.Fprint(w, "FAILURE\t\n")
			FatalFlush(w, err)
		}
		if !hasIntermediate {
			if err = vaulter.Mount(vaultAPI, mount, &vaulter.MountConfiguration{
				Type:        "pki",
				Description: "testing intermediate CA",
				MaxLeaseTTL: "26280h",
			}); err != nil {
				fmt.Fprint(w, "FAILURE\t\n")
				FatalFlush(w, err)
			}
		}
		fmt.Fprint(w, "SUCCESS\t\n")

		fmt.Fprintf(w, "Creating a CSR:\t")
		csrConfig := &vaulter.CSRConfig{
			CommonName: commonName,
			TTL:        "26280h",
			KeyBits:    4096,
		}
		csrSecret, err := vaulter.CSR(vaultAPI, mount, csrConfig)
		if err != nil {
			fmt.Fprint(w, "FAILURE\t\n")
			FatalFlush(w, err)
		}
		fmt.Fprint(w, "SUCCESS\t\n")

		fmt.Fprint(w, "Signing the intermediate CSR with the root CA:\t")
		csr := csrSecret.Data["csr"].(string)
		csrSigningConfig := &vaulter.CSRSigningConfig{
			CommonName: commonName,
			TTL:        "8760h",
		}
		fmt.Println(rootmount)
		signedCert, err := vaulter.SignCSR(vaultAPI, rootmount, csr, csrSigningConfig)
		if err != nil {
			fmt.Fprint(w, "FAILURE\t\n")
			FatalFlush(w, err)
		}
		fmt.Fprint(w, "SUCCESS\t\n")

		fmt.Fprint(w, "Importing the signed cert into the intermediate CA:\t")
		certContents := signedCert.Data["certificate"].(string)
		_, err = vaulter.ImportCert(vaultAPI, mount, certContents)
		if err != nil {
			fmt.Fprint(w, "FAILURE\t\n")
			FatalFlush(w, err)
		}
		fmt.Fprint(w, "SUCCESS\t\n")

		fmt.Fprint(w, "Set the CA and CRL URLs for the intermediate CA:\t")
		urlParts, err := url.Parse(vaultURL)
		if err != nil {
			FatalFlush(w, err)
		}
		_, err = vaulter.ConfigCAAccess(
			vaultAPI,
			urlParts.Scheme,
			fmt.Sprintf("%s:%s", urlParts.Hostname(), urlParts.Port()),
			mount,
		)
		if err != nil {
			fmt.Fprint(w, "FAILURE\t\n")
			FatalFlush(w, err)
		}
		fmt.Fprint(w, "SUCCESS\t\n")
		w.Flush()
	},
}

var intermediateCACheckCmd = &cobra.Command{
	Use:   "intermediate-ca",
	Short: "Checks the status of the intermediate CA in Vault",
	Long: `Checks the status of the intermediate CA in Vault by determining the
following:
    1. If the intermediate CA backend is mounted.
    2. If the role exists.
    3. If the intermediate CA backend is configured correctly.
This command does not create any of the above if it does not exist. If the
backend is not mounted, then the status of each subsequent check will be
'UNKNOWN'.`,
	Run: func(cmd *cobra.Command, args []string) {
		if mount == "" {
			log.Fatal("--mount was not set.")
		}
		if role == "" {
			log.Fatal("--role was not set.")
		}
		if commonName == "" {
			log.Fatal("--common-name was not set.")
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)

		fmt.Fprint(w, "Intermediate CA backend is mounted:\t")
		hasIntermediate, err := vaulter.IsMounted(vaultAPI, mount)
		if err != nil {
			FatalFlush(w, err)
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
				FatalFlush(w, err)
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
				FatalFlush(w, err)
			}
			if config == nil {
				FatalFlush(w, errors.New("config was nil"))
			}
			var (
				v  string
				ok bool
			)
			if v, ok = config["issuing_certificates"].(string); !ok {
				FatalFlush(w, errors.New("issuing_certificates was not found"))
			}
			if v != fmt.Sprintf("%s/v1/%s/ca", vaultURL, mount) {
				fmt.Fprint(w, "FAILURE\t\n")
				FatalFlush(w, fmt.Errorf("issuing_certificates was %s", v))
			}

			if v, ok = config["crl_distribution_points"].(string); !ok {
				fmt.Fprint(w, "FAILURE\t\n")
				FatalFlush(w, errors.New("crl_distribution_points was not found"))
			}
			if v != fmt.Sprintf("%s/v1/%s/crl", vaultURL, mount) {
				fmt.Fprintf(w, "FAILURE\t\n")
				FatalFlush(w, fmt.Errorf("crl_distribution_points was %s", v))
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
		if mount == "" {
			log.Fatal("--mount must be set.")
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)

		fmt.Fprint(w, "Unmounting intermediate CA backend:\t")
		hasRoot, err := vaulter.IsMounted(vaultAPI, mount)
		if err != nil {
			fmt.Fprintf(w, "FAILURE\t\n")
			FatalFlush(w, err)
		}
		if hasRoot {
			if err = vaulter.Unmount(vaultAPI, mount); err != nil {
				fmt.Fprintf(w, "FAILURE\t\n")
				FatalFlush(w, err)
			}
			fmt.Fprintf(w, "SUCCESS\t\n")
		} else {
			fmt.Fprintf(w, "SUCCESS\t\n")
		}
		w.Flush()
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
		defaultIntMount,
		"The path in Vault to the intermediate CA pki backend.",
	)
	intermediateCAInitCmd.PersistentFlags().StringVar(
		&rootmount,
		"root-mount",
		defaultRootMount,
		"The paht in Vault to the root CA pki backend.",
	)
	intermediateCAInitCmd.PersistentFlags().StringVar(
		&role, // defined in root.go
		"role",
		defaultIntRole,
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
		defaultIntMount,
		"The path in Vault to the intermediate CA pki backend.",
	)
	intermediateCACheckCmd.PersistentFlags().StringVar(
		&role, // defined in root.go
		"role",
		defaultIntRole,
		"The name of the role to use for operations on the intermediate CA.",
	)
	intermediateCACheckCmd.PersistentFlags().StringVar(
		&commonName, // defined in root.go
		"common-name",
		"",
		"The common name to use for operations on the intermediate CA.",
	)
}
