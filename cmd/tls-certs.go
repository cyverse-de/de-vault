package cmd

import "github.com/spf13/cobra"

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

}

func (t *TLSGen) revokeRun(cmd *cobra.Command, args []string) {

}

func init() {
	t := NewTLSGen()
	generateCmd.AddCommand(t.Generate)
	checkCmd.AddCommand(t.Check)
	revokeCmd.AddCommand(t.Revoke)
}
