// Package commands holds the different CLI commands
package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/bytemare/goproject/internal/version"

	"github.com/bytemare/goproject/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// profileCmd represents the profile command
func profileCmd() *cobra.Command {
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage your profiles",
		Long: `The profile command lets you manage profiles that are used by the new command.
It allows you to create, edit and list profiles, and set the default profile to use.`,
		Args: cobra.MinimumNArgs(1),
	}

	profileCmd.AddCommand(profileCmdNew())
	profileCmd.AddCommand(profileCmdList())
	profileCmd.AddCommand(profileCmdEdit())
	profileCmd.AddCommand(profileCmdDefault())

	return profileCmd
}

func profileCmdNew() *cobra.Command {
	return &cobra.Command{
		Use:   "new [new profile name]",
		Short: "create a new profile",
		Long: `new creates a new profile, opening an editor to configure it.
The command returns immediately with an error if a profile with the same name already exists.`,
		Args: cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			// Update
			if viper.Get(config.DefaultAutoUpdateKeyName) == true {
				version.Upgrade()
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			profileNew(args[0])
		},
	}
}

func profileCmdList() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list registered profiles",
		Long:  `list registered profiles as found in the profiles directory`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			profileList()
		},
	}
}

func profileCmdEdit() *cobra.Command {
	return &cobra.Command{
		Use:     "edit",
		Short:   "edit existing profile",
		Long:    "edit existing profile as found in the profiles directory",
		Example: "edit myprofile",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			profileEdit(args[0])
		},
	}
}

func profileCmdDefault() *cobra.Command {
	return &cobra.Command{
		Use:   "default [existing profile name]",
		Short: "sets default profile",
		Long:  `default sets the specified profile as the default profile to use`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			profileSetDefault(args[0])
		},
	}
}

func profileList() {
	profiles, err := config.ListProfiles()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(profiles) == 0 {
		fmt.Println("No profiles are registered.")
		os.Exit(0)
	}

	fmt.Println("Registered profiles :")

	defaultProfile := viper.GetString(config.DefaultConfigProfileKeyName)
	defaultProfile = strings.TrimSuffix(defaultProfile, ".toml")

	for _, p := range profiles {
		if p == defaultProfile {
			fmt.Printf("\t- %s [default]\n", p)
		} else {
			fmt.Printf("\t- %s\n", p)
		}
	}

	os.Exit(0)
}

func profileSetDefault(name string) {
	err := config.SetDefaultProfile(name)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Default profile set to %s.\n", name)

	os.Exit(0)
}

func profileNew(name string) {
	filepath, err := config.CreateNewProfileFile(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err = config.EditProfile(name); err != nil {
		fmt.Println(err)

		if err = os.Remove(filepath); err != nil {
			fmt.Printf("Could not remove %s : %s\n", filepath, err)
		}

		os.Exit(1)
	}

	fmt.Printf("Profile %s was successfully created.\n", name)
	os.Exit(0)
}

func profileEdit(name string) {
	if err := config.EditProfile(name); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
