package cmd

import (
	"log"
	"net/url"

	"github.com/cyverse-de/de-vault/vaulter"
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

const defaultRootRole = "root-ca"
const defaultRootMount = "root-ca"
const defaultIntRole = "intermediate-ca"
const defaultIntMount = "intermediate-ca"

// Flusher can flush stuff.
type Flusher interface {
	Flush() error
}

// FatalFlush flushs something and the exits with a log.Fatal() call.
func FatalFlush(f Flusher, e error) {
	f.Flush()
	log.Fatal(e)
}

// RootCmd is the root node in the command tree
var RootCmd = &cobra.Command{
	Use:   "de-vault",
	Short: "Utility for managing Vault for the Discovery Environment",
	Long: `A command-line utility for managing a deployment of Hashicorp's Vault
project. This tool is geared towards CyVerse's Discovery Environment.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if parentToken == "" {
			log.Fatal("--token must be set.")
		}

		if vaultURL == "" {
			log.Fatal("--api-url must be set.")
		}
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
	certPath string // Writable path to a file that will contain a TLS cert.
	keyPath  string // Writable path to a file that will contain a TLS key.
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&parentToken, "token", "t", "", "The Vault parent token.")
	RootCmd.PersistentFlags().StringVarP(&vaultURL, "api-url", "u", "http://127.0.0.1:8200", "The URL for the Vault API.")
	RootCmd.PersistentFlags().StringVarP(&clientCert, "client-cert", "c", "", "The client TLS certificate to use for the Vault connection.")
	RootCmd.PersistentFlags().StringVarP(&clientKey, "client-key", "k", "", "The client key to use for TLS connection to the Vault API.")
}
