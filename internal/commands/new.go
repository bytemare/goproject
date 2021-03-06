// Package commands holds the different CLI commands
package commands

import (
	"fmt"
	"os"

	"github.com/bytemare/goproject/internal/config"
	"github.com/bytemare/goproject/internal/templates"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// newCmd represents the new command
func newCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new go project",
		Long: `This command will create a new project in the current directory,
given the layout of your profile.
For example :

./goproject new myApp

Will load the default profile as specified in your configuration settings,
and create the files and directories as specified in the profile.

./goproject new myApp -p myProfile

Will do the same but with the specified profile
`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			setupNewProject(cmd, args)

			if err := checkVars(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			newProject()
		},
	}

	newCmd.Flags().StringP("profile", "p", "", "Specify the profile you want to use")

	return newCmd
}

func setupNewProject(cmd *cobra.Command, args []string) {
	// If no argument was given, we develop the project inside the current directory, thus inheriting its name
	if len(args) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Unable to get working directory : %s\n", err)
			os.Exit(1)
		}

		viper.Set("name", wd)
	} else {
		viper.Set("name", args[0])
	}

	if cmd.Flag("profile").Value.String() != "" {
		viper.Set("profile", cmd.Flag("profile").Value.String())
	}
}

// checkVars verifies all necessary information is given to build a new project
func checkVars() error {
	// todo project name, profile
	return nil
}

func newProject() {
	// If profile is given
	profileName := viper.GetString("profile")
	if profileName == "" {
		fmt.Println("Loading default profile")
		// No profile was specified, we're therefore calling the default profile
		profileName = viper.GetString(config.DefaultConfigProfileKeyName)
		if profileName == "" {
			fmt.Println("Error : no profile was specified, and no default profile was found.")
			os.Exit(1)
		}
	}

	prof, err := config.LoadProfile(profileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initiate and create project
	projectName := viper.GetString("name")

	projectLocation := viper.GetString("location")
	if projectLocation == "" {
		projectLocation = config.DefaultTargetProjectLocation
	}

	project := templates.NewProject(prof, projectName, projectLocation)

	// Build project
	if err := project.Build(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Project %s was successfully created.\n", projectName)
	os.Exit(0)
}
