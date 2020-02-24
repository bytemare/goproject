package cmd

import (
	"github.com/bytemare/goproject/internal"
	"github.com/spf13/cobra"
)

func upgradeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade",
		Short: "upgrade goproject tool to latest version",
		Long:  "upgrade goproject tool to latest version",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			internal.Upgrade()
		},
	}
}
