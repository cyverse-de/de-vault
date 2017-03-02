package cmd

import (
	"log"
	"net/url"

	"github.com/cyverse-de/vaulter"
	"github.com/spf13/cobra"
)

var (
	parentToken string
	vaultURL    string
	clientCert  string
	clientKey   string
	vaultAPI    *vaulter.VaultAPI
	vaultCFG    *vaulter.VaultAPIConfig
)

// RootCmd is the root node in the command tree
var RootCmd = &cobra.Command{
	Use:   "de-vault",
	Short: "Utility for managing Vault for the Discovery Environment",
	Long: `A command-line utility for managing a deployment of Hashicorp's Vault
project. This tool is geared towards CyVerse's Discovery Environment.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		connURL, err := url.Parse(vaultURL)
		if err != nil {
			log.Fatal(err)
		}
		vaultCFG = &vaulter.VaultAPIConfig{
			ParentToken: parentToken,
			Host:        connURL.Hostname(),
			Port:        connURL.Port(),
			Scheme:      connURL.Scheme,
			ClientCert:  clientCert,
			ClientKey:   clientKey,
		}
		vaultAPI = &vaulter.VaultAPI{}
		if err = vaulter.InitAPI(vaultAPI, vaultCFG, vaultCFG.ParentToken); err != nil {
			log.Fatal(err)
		}
	},
}

var (
	mount      string // Path to a backend in Vault.
	role       string // Name of the role used in some operations in Vault.
	commonName string // The CN to use for some TLS-related operations.
	certPath   string // Writable path to a file that will contain a TLS cert.
	keyPath    string // Writable path to a file that will contain a TLS key.
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&parentToken, "token", "t", "", "The Vault parent token.")
	RootCmd.PersistentFlags().StringVarP(&vaultURL, "api-url", "u", "http://127.0.0.1:8200", "The URL for the Vault API.")
	RootCmd.PersistentFlags().StringVarP(&clientCert, "client-cert", "c", "", "The client TLS certificate to use for the Vault connection.")
	RootCmd.PersistentFlags().StringVarP(&clientKey, "client-key", "k", "", "The client key to use for TLS connection to the Vault API.")
}
