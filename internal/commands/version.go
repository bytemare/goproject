// Package commands holds the different CLI commands
package commands

import (
	"fmt"

	"github.com/bytemare/goproject/internal/version"

	"github.com/spf13/cobra"
)

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "displays the actual version",
		Long:  "version will display the version and commit the command was built on",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.PrintableVersion())
		},
	}
}
