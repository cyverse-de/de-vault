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

// IntermediateCA contains the functionality associated with checking,
// initializing, and removing intermediate CAs.
type IntermediateCA struct {
	rootMount  string
	mount      string
	role       string
	commonName string
	Init       *cobra.Command
	Check      *cobra.Command
	Remove     *cobra.Command
}

// NewIntermediateCA returns a newly initialized *IntermediateCA.
func NewIntermediateCA() *IntermediateCA {
	ca := &IntermediateCA{
		Init: &cobra.Command{
			Use:   "intermediate-ca",
			Short: "Initialize an intermediate CA in Vault.",
			Long: `Initializes an intermediate CA in Vault, the end result being a
new PKI backend that has a role configured and a signed CSR imported into it.`,
		},
		Check: &cobra.Command{
			Use:   "intermediate-ca",
			Short: "Checks the status of the intermediate CA in Vault",
			Long: `Checks the status of the intermediate CA in Vault by determining
			the following:
				1. If the intermediate CA backend is mounted.
				2. If the role exists.
				3. If the intermediate CA backend is configured correctly.
			This command does not create any of the above if it does not exist. If the
			backend is not mounted, then the status of each subsequent check will be
			'UNKNOWN'.`,
		},
		Remove: &cobra.Command{
			Use:   "intermediate-ca",
			Short: "A brief description of your command",
			Long: `A longer description that spans multiple lines and likely contains
			examples and usage of using your command. For example:
			Cobra is a CLI library for Go that empowers applications. This application
			is a tool to generate the needed files to quickly create a Cobra
			application.`,
		},
	}
	ca.Init.Run = ca.initRun
	ca.Check.Run = ca.checkRun
	ca.Remove.Run = ca.removeRun

	ca.Init.PersistentFlags().StringVar(
		&ca.mount, // defined in root.go
		"mount",
		defaultIntMount,
		"The path in Vault to the intermediate CA pki backend.",
	)
	ca.Init.PersistentFlags().StringVar(
		&ca.rootMount,
		"root-mount",
		defaultRootMount,
		"The path in Vault to the root CA pki backend.",
	)
	ca.Init.PersistentFlags().StringVar(
		&ca.role, // defined in root.go
		"role",
		defaultIntRole,
		"The name of the role to use for operations on the intermediate CA.",
	)
	ca.Init.PersistentFlags().StringVar(
		&ca.commonName, // defined in root.go
		"common-name",
		"",
		"The common name to use for operations on the intermediate CA.",
	)

	ca.Check.PersistentFlags().StringVar(
		&ca.mount, // defined in root.go
		"mount",
		defaultIntMount,
		"The path in Vault to the intermediate CA pki backend.",
	)
	ca.Check.PersistentFlags().StringVar(
		&ca.role, // defined in root.go
		"role",
		defaultIntRole,
		"The name of the role to use for operations on the intermediate CA.",
	)
	ca.Check.PersistentFlags().StringVar(
		&ca.commonName, // defined in root.go
		"common-name",
		"",
		"The common name to use for operations on the intermediate CA.",
	)

	ca.Remove.PersistentFlags().StringVar(
		&ca.mount, // defined in root.go
		"mount",
		defaultIntMount,
		"The path in Vault to the intermediate CA pki backend.",
	)
	ca.Remove.PersistentFlags().StringVar(
		&ca.role, // defined in root.go
		"role",
		defaultIntRole,
		"The name of the role to use for operations on the intermediate CA.",
	)
	ca.Remove.PersistentFlags().StringVar(
		&ca.commonName, // defined in root.go
		"common-name",
		"",
		"The common name to use for operations on the intermediate CA.",
	)

	return ca
}

func (i *IntermediateCA) initRun(cmd *cobra.Command, args []string) {
	if i.mount == "" {
		log.Fatal("--mount was not set.")
	}
	if i.role == "" {
		log.Fatal("--role was not set.")
	}
	if i.commonName == "" {
		log.Fatal("--common-name was not set.")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)

	fmt.Fprintf(w, "Creating the intermediate CA:\t")
	hasIntermediate, err := vaulter.IsMounted(vaultAPI, i.mount)
	if err != nil {
		fmt.Fprint(w, "FAILURE\t\n")
		FatalFlush(w, err)
	}
	if !hasIntermediate {
		if err = vaulter.Mount(vaultAPI, i.mount, &vaulter.MountConfiguration{
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
		CommonName: i.commonName,
		TTL:        "26280h",
		KeyBits:    4096,
	}
	csrSecret, err := vaulter.CSR(vaultAPI, i.mount, csrConfig)
	if err != nil {
		fmt.Fprint(w, "FAILURE\t\n")
		FatalFlush(w, err)
	}
	fmt.Fprint(w, "SUCCESS\t\n")

	fmt.Fprint(w, "Signing the intermediate CSR with the root CA:\t")
	csr := csrSecret.Data["csr"].(string)
	csrSigningConfig := &vaulter.CSRSigningConfig{
		CommonName: i.commonName,
		TTL:        "8760h",
	}
	signedCert, err := vaulter.SignCSR(vaultAPI, i.rootMount, csr, csrSigningConfig)
	if err != nil {
		fmt.Fprint(w, "FAILURE\t\n")
		FatalFlush(w, err)
	}
	fmt.Fprint(w, "SUCCESS\t\n")

	fmt.Fprint(w, "Importing the signed cert into the intermediate CA:\t")
	certContents := signedCert.Data["certificate"].(string)
	_, err = vaulter.ImportCert(vaultAPI, i.mount, certContents)
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
		i.mount,
	)
	if err != nil {
		fmt.Fprint(w, "FAILURE\t\n")
		FatalFlush(w, err)
	}
	fmt.Fprint(w, "SUCCESS\t\n")
	w.Flush()
}

func (i *IntermediateCA) checkRun(cmd *cobra.Command, args []string) {
	if i.mount == "" {
		log.Fatal("--mount was not set.")
	}
	if i.role == "" {
		log.Fatal("--role was not set.")
	}
	if i.commonName == "" {
		log.Fatal("--common-name was not set.")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)

	fmt.Fprint(w, "Intermediate CA backend is mounted:\t")
	hasIntermediate, err := vaulter.IsMounted(vaultAPI, i.mount)
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
		hasRole, err := vaulter.HasRole(vaultAPI, i.mount, i.role, i.commonName, true)
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
		configSecret, err := vaultAPI.Read(vaultAPI.Client(), fmt.Sprintf("%s/config/urls", i.mount))
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
		expected := fmt.Sprintf("%s/v1/%s/ca", vaultURL, i.mount)
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
		expected = fmt.Sprintf("%s/v1/%s/crl", vaultURL, i.mount)
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
}

func (i *IntermediateCA) removeRun(cmd *cobra.Command, args []string) {
	if i.mount == "" {
		log.Fatal("--mount must be set.")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)

	fmt.Fprint(w, "Unmounting intermediate CA backend:\t")
	hasRoot, err := vaulter.IsMounted(vaultAPI, i.mount)
	if err != nil {
		fmt.Fprintf(w, "FAILURE\t\n")
		FatalFlush(w, err)
	}
	if hasRoot {
		if err = vaulter.Unmount(vaultAPI, i.mount); err != nil {
			fmt.Fprintf(w, "FAILURE\t\n")
			FatalFlush(w, err)
		}
		fmt.Fprintf(w, "SUCCESS\t\n")
	} else {
		fmt.Fprintf(w, "SUCCESS\t\n")
	}
	w.Flush()
}

var intermediate *IntermediateCA

func init() {
	intermediate = NewIntermediateCA()
	removeCmd.AddCommand(intermediate.Remove)
	initCmd.AddCommand(intermediate.Init)
	checkCmd.AddCommand(intermediate.Check)
}
