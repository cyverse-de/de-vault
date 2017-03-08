package cmd

import "github.com/spf13/cobra"

// TLSGen contains the commands for managing TLS certs and keys.
type TLSGen struct {
	mount      string
	role       string
	commonName string
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
