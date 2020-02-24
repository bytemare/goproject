package cmd

import (
	"fmt"
	"os"

	"github.com/bytemare/goproject/internal/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
func getRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "goproject",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Load and initialise configuration
			config.Initialise()
		},
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd := getRootCmd()
	rootCmd.AddCommand(newCmd())
	rootCmd.AddCommand(profileCmd())
	rootCmd.AddCommand(versionCmd())
	rootCmd.AddCommand(upgradeCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
