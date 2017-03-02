package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// tls-certsCmd represents the tls-certs command
var tlsCertsCmd = &cobra.Command{
	Use:   "tls-certs",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("tls-certs called")
	},
}

func init() {
	generateCmd.AddCommand(tlsCertsCmd)
	tlsCertsCmd.PersistentFlags().StringVar(
		&certPath,
		"tls-cert",
		"",
		"The path that will contain the new TLS cert.",
	)
	tlsCertsCmd.PersistentFlags().StringVar(
		&keyPath,
		"tls-key",
		"",
		"The path that will contain the new TLS key.",
	)
}
