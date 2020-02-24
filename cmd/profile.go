package cmd

import (
	"fmt"
	"os"

	"github.com/bytemare/goproject/internal"
	"github.com/spf13/viper"

	"github.com/bytemare/goproject/internal/config"
	"github.com/spf13/cobra"
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
				internal.Upgrade()
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

/*
func init() {
	rootCmd.AddCommand(profileCmd)

	profileCmd.AddCommand(profileCmdNew)
	profileCmd.AddCommand(profileCmdList)
	profileCmd.AddCommand(profileCmdDefault)
}*/

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
	for _, p := range profiles {
		fmt.Printf("\t- %s\n", p)
	}

	os.Exit(0)
}

func profileSetDefault(name string) {
	err := config.SetDefaultProfile(name)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func profileNew(name string) {
	filepath, err := config.CreateNewProfileFile(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err = config.NanoProfile(name); err != nil {
		_ = os.Remove(filepath)
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Profile %s was successfully created.\n", name)
	os.Exit(0)
}
