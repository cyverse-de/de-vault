package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// root-caCmd represents the root-ca command
var rootCaCmd = &cobra.Command{
	Use:   "root-ca",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("root-ca called")
	},
}

func init() {
	initCmd.AddCommand(rootCaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// root-caCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// root-caCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
