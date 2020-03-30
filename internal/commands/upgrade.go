// Package commands holds the different CLI commands
package commands

import (
	"github.com/bytemare/goproject/internal/version"

	"github.com/spf13/cobra"
)

func upgradeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade",
		Short: "upgrade goproject tool to latest version",
		Long:  "upgrade goproject tool to latest version",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			version.Upgrade()
		},
	}
}
