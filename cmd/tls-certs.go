package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"text/tabwriter"

	"github.com/cyverse-de/vaulter"
	"github.com/spf13/cobra"
)

// TLSGen contains the commands for managing TLS certs and keys.
type TLSGen struct {
	mount      string
	role       string
	commonName string
	certPath   string
	keyPath    string
	Check      *cobra.Command
	Generate   *cobra.Command
	Revoke     *cobra.Command
}

// NewTLSGen returns a newly instantiated *TLSGen.
func NewTLSGen() *TLSGen {
	t := &TLSGen{
		Check: &cobra.Command{
			Use:   "tls",
			Short: "Checks the status of a TLS cert/key pair by the serial number.",
			Long:  "Checks the status of a TLS cert/key pair by the serial number.",
		},
		Generate: &cobra.Command{
			Use:   "tls",
			Short: "Generate a new TLS cert/key pair.",
			Long:  "Generates a new TLS cert/key pair.",
		},
		Revoke: &cobra.Command{
			Use:   "tls",
			Short: "Revokes a TLS cert/key pair.",
			Long:  "Revokes a TLS cert/key pair.",
		},
	}

	t.Check.Run = t.checkRun
	t.Generate.Run = t.generateRun
	t.Revoke.Run = t.revokeRun

	t.Generate.PersistentFlags().StringVar(
		&t.mount,
		"mount",
		defaultIntMount,
		"The path in Vault to the intermediate CA backend.",
	)
	t.Generate.PersistentFlags().StringVar(
		&t.role,
		"role",
		"",
		"The role to create for generating TLS certs/keys. Should be different for each site.",
	)
	t.Generate.PersistentFlags().StringVar(
		&t.commonName,
		"common-name",
		"",
		"The common name to use when generating the TLS certs/keys.",
	)
	t.Generate.PersistentFlags().StringVar(
		&t.certPath,
		"cert-path",
		"",
		"The file path for the TLS cert. Should be writable.",
	)
	t.Generate.PersistentFlags().StringVar(
		&t.keyPath,
		"key-path",
		"",
		"The file path for the TLS key. Should be writable.",
	)

	return t
}

func (t *TLSGen) checkRun(cmd *cobra.Command, args []string) {

}

func (t *TLSGen) generateRun(cmd *cobra.Command, args []string) {
	var err error
	if t.mount == "" {
		log.Fatal("--mount must be set.")
	}
	if t.role == "" {
		log.Fatal("--role must be set.")
	}
	if t.commonName == "" {
		log.Fatal("--common-name must be set.")
	}
	if t.certPath == "" {
		log.Fatal("--cert-path must be set.")
	}
	if t.keyPath == "" {
		log.Fatal("--key-path must be set.")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)

	fmt.Fprint(w, "Create a role for cert generation: \t")
	certRoleConfig := &vaulter.RoleConfig{
		KeyBits:      4096,
		MaxTTL:       "8760h",
		AllowAnyName: true,
	}
	if _, err = vaulter.CreateRole(vaultAPI, t.mount, t.role, certRoleConfig); err != nil {
		fmt.Fprint(w, "FAILURE\t\n")
		FatalFlush(w, err)
	}
	fmt.Fprint(w, "SUCCESS\t\n")

	fmt.Fprint(w, "Create a cert with the role:\t")
	issueCertConfig := &vaulter.IssueCertConfig{
		CommonName: t.commonName,
		TTL:        "720h",
		Format:     "pem",
	}
	certSecret, err := vaulter.IssueCert(vaultAPI, t.mount, t.role, issueCertConfig)
	if err != nil {
		fmt.Fprint(w, "FAILURE\t\n")
		FatalFlush(w, err)
	}

	if _, ok := certSecret.Data["certificate"]; !ok {
		fmt.Fprint(w, "FAILURE\t\n")
		FatalFlush(w, errors.New("No certificate found."))
	}

	if _, ok := certSecret.Data["issuing_ca"]; !ok {
		fmt.Fprint(w, "FAILURE\t\n")
		FatalFlush(w, errors.New("No issuing CA found."))
	}
	if _, ok := certSecret.Data["serial_number"]; !ok {
		fmt.Fprint(w, "FAILURE\t\n")
		FatalFlush(w, errors.New("No serial number found."))
	}
	fmt.Fprint(w, "SUCCESS\t\n")

	fmt.Fprint(w, "Writing cert to file:\t")
	certfile, err := os.Create(t.certPath)
	defer certfile.Close()
	if err != nil {
		fmt.Fprint(w, "FAILURE\t\n")
		FatalFlush(w, err)
	}
	if _, err = io.WriteString(certfile, certSecret.Data["certificate"].(string)); err != nil {
		fmt.Fprint(w, "FAILURE\t\n")
		FatalFlush(w, err)
	}
	if _, err = io.WriteString(certfile, "\n"); err != nil {
		fmt.Fprint(w, "FAILURE\t\n")
		FatalFlush(w, err)
	}
	if _, err = io.WriteString(certfile, certSecret.Data["issuing_ca"].(string)); err != nil {
		fmt.Fprint(w, "FAILURE\t\n")
		FatalFlush(w, err)
	}
	if _, err = io.WriteString(certfile, "\n"); err != nil {
		fmt.Fprint(w, "FAILURE\t\n")
		FatalFlush(w, err)
	}
	fmt.Fprint(w, "SUCCESS\t\n")

	fmt.Fprint(w, "Write key to file:\t")
	if _, ok := certSecret.Data["private_key"]; !ok {
		fmt.Fprint(w, "FAILURE\t\n")
		FatalFlush(w, err)
	}

	keyfile, err := os.Create(t.keyPath)
	defer keyfile.Close()
	if err != nil {
		fmt.Fprint(w, "FAILURE\t\n")
		FatalFlush(w, err)
	}
	if _, err = io.WriteString(keyfile, certSecret.Data["private_key"].(string)); err != nil {
		FatalFlush(w, err)
	}
	if _, err = io.WriteString(keyfile, "\n"); err != nil {
		fmt.Fprint(w, "FAILURE\t\n")
		FatalFlush(w, err)
	}
	fmt.Fprint(w, "SUCCESS\t\n")

	fmt.Fprint(w, "TLS cert/key serial number (SAVE THIS):\t")
	fmt.Fprint(w, fmt.Sprintf("%s\t\n", certSecret.Data["serial_number"]))

	w.Flush()
}

func (t *TLSGen) revokeRun(cmd *cobra.Command, args []string) {

}

func init() {
	t := NewTLSGen()
	generateCmd.AddCommand(t.Generate)
	checkCmd.AddCommand(t.Check)
	revokeCmd.AddCommand(t.Revoke)
}
