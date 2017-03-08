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
	intRootMount  string
	intMount      string
	intRole       string
	intCommonName string
)

var intermediateCAInitCmd = &cobra.Command{
	Use:   "intermediate-ca",
	Short: "Initialize an intermediate CA in Vault.",
	Long: `Initializes an intermediate CA in Vault, the end result being a new PKI
 backend that has a role configured and a signed CSR imported into it.`,
	Run: func(cmd *cobra.Command, args []string) {
		if intMount == "" {
			log.Fatal("--mount was not set.")
		}
		if intRole == "" {
			log.Fatal("--role was not set.")
		}
		if intCommonName == "" {
			log.Fatal("--common-name was not set.")
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)

		fmt.Fprintf(w, "Creating the intermediate CA:\t")
		hasIntermediate, err := vaulter.IsMounted(vaultAPI, intMount)
		if err != nil {
			fmt.Fprint(w, "FAILURE\t\n")
			FatalFlush(w, err)
		}
		if !hasIntermediate {
			if err = vaulter.Mount(vaultAPI, intMount, &vaulter.MountConfiguration{
				Type:        "pki",
				Description: "intermediate CA",
				MaxLeaseTTL: "26280h",
			}); err != nil {
				fmt.Fprint(w, "FAILURE\t\n")
				FatalFlush(w, err)
			}
		}
		fmt.Fprint(w, "SUCCESS\t\n")

		fmt.Fprintf(w, "Creating a CSR:\t")
		csrConfig := &vaulter.CSRConfig{
			CommonName: intCommonName,
			TTL:        "26280h",
			KeyBits:    4096,
		}
		csrSecret, err := vaulter.CSR(vaultAPI, intMount, csrConfig)
		if err != nil {
			fmt.Fprint(w, "FAILURE\t\n")
			FatalFlush(w, err)
		}
		fmt.Fprint(w, "SUCCESS\t\n")

		fmt.Fprint(w, "Signing the intermediate CSR with the root CA:\t")
		csr := csrSecret.Data["csr"].(string)
		csrSigningConfig := &vaulter.CSRSigningConfig{
			CommonName: intCommonName,
			TTL:        "8760h",
		}
		fmt.Println(intRootMount)
		signedCert, err := vaulter.SignCSR(vaultAPI, intRootMount, csr, csrSigningConfig)
		if err != nil {
			fmt.Fprint(w, "FAILURE\t\n")
			FatalFlush(w, err)
		}
		fmt.Fprint(w, "SUCCESS\t\n")

		fmt.Fprint(w, "Importing the signed cert into the intermediate CA:\t")
		certContents := signedCert.Data["certificate"].(string)
		_, err = vaulter.ImportCert(vaultAPI, intMount, certContents)
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
			intMount,
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
		if intMount == "" {
			log.Fatal("--mount was not set.")
		}
		if intRole == "" {
			log.Fatal("--role was not set.")
		}
		if intCommonName == "" {
			log.Fatal("--common-name was not set.")
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)

		fmt.Fprint(w, "Intermediate CA backend is mounted:\t")
		hasIntermediate, err := vaulter.IsMounted(vaultAPI, intMount)
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
			hasRole, err := vaulter.HasRole(vaultAPI, intMount, intRole, intCommonName, true)
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
			configSecret, err := vaultAPI.Read(vaultAPI.Client(), fmt.Sprintf("%s/config/urls", intMount))
			if err != nil {
				FatalFlush(w, err)
			}
			if configSecret == nil {
				FatalFlush(w, errors.New("config was nil"))
			}
			if configSecret.Data == nil {
				FatalFlush(w, errors.New("config.Data was nil"))
			}
			var (
				ok            bool
				issuingCerts  []interface{}
				crlDistPoints []interface{}
				config        map[string]interface{}
			)
			config = configSecret.Data
			if issuingCerts, ok = config["issuing_certificates"].([]interface{}); !ok {
				FatalFlush(w, errors.New("issuing_certificates was not found"))
			}
			expected := fmt.Sprintf("%s/v1/%s/ca", vaultURL, intMount)
			found := false
			for _, c := range issuingCerts {
				found = found || c.(string) == expected
			}
			if !found {
				fmt.Fprint(w, "NO\t\n")
				FatalFlush(w, errors.New("issuing_certificates was %s"))
			}
			if crlDistPoints, ok = config["crl_distribution_points"].([]interface{}); !ok {
				fmt.Fprint(w, "NO\t\n")
				FatalFlush(w, errors.New("crl_distribution_points was not found"))
			}
			expected = fmt.Sprintf("%s/v1/%s/crl", vaultURL, intMount)
			found = false
			for _, d := range crlDistPoints {
				found = found || d.(string) == expected
			}
			if !found {
				fmt.Fprintf(w, "NO\t\n")
				FatalFlush(w, errors.New("crl_distribution_points was not found"))
			}
			fmt.Fprintf(w, "YES\t\n")
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
		if intMount == "" {
			log.Fatal("--mount must be set.")
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)

		fmt.Fprint(w, "Unmounting intermediate CA backend:\t")
		hasRoot, err := vaulter.IsMounted(vaultAPI, intMount)
		if err != nil {
			fmt.Fprintf(w, "FAILURE\t\n")
			FatalFlush(w, err)
		}
		if hasRoot {
			if err = vaulter.Unmount(vaultAPI, intMount); err != nil {
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
		&intMount, // defined in root.go
		"mount",
		defaultIntMount,
		"The path in Vault to the intermediate CA pki backend.",
	)
	intermediateCARemoveCmd.PersistentFlags().StringVar(
		&intRole, // defined in root.go
		"role",
		defaultIntRole,
		"The name of the role to use for operations on the intermediate CA.",
	)
	intermediateCARemoveCmd.PersistentFlags().StringVar(
		&intCommonName, // defined in root.go
		"common-name",
		"",
		"The common name to use for operations on the intermediate CA.",
	)

	initCmd.AddCommand(intermediateCAInitCmd)
	intermediateCAInitCmd.PersistentFlags().StringVar(
		&intMount, // defined in root.go
		"mount",
		defaultIntMount,
		"The path in Vault to the intermediate CA pki backend.",
	)
	intermediateCAInitCmd.PersistentFlags().StringVar(
		&intRootMount,
		"root-mount",
		defaultRootMount,
		"The path in Vault to the root CA pki backend.",
	)
	intermediateCAInitCmd.PersistentFlags().StringVar(
		&intRole, // defined in root.go
		"role",
		defaultIntRole,
		"The name of the role to use for operations on the intermediate CA.",
	)
	intermediateCAInitCmd.PersistentFlags().StringVar(
		&intCommonName, // defined in root.go
		"common-name",
		"",
		"The common name to use for operations on the intermediate CA.",
	)

	intermediateCACheckCmd.Flags().StringVar(
		&intMount, // defined in root.go
		"mount",
		defaultIntMount,
		"The path in Vault to the intermediate CA pki backend.",
	)
	intermediateCACheckCmd.Flags().StringVar(
		&intRole, // defined in root.go
		"role",
		defaultIntRole,
		"The name of the role to use for operations on the intermediate CA.",
	)
	intermediateCACheckCmd.Flags().StringVar(
		&intCommonName, // defined in root.go
		"common-name",
		"",
		"The common name to use for operations on the intermediate CA.",
	)
	checkCmd.AddCommand(intermediateCACheckCmd)
}
