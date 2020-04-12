// Package commands holds the different CLI commands
package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"

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
			name, location, profile, err := setupNewProject(cmd, args)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			newProject(name, location, profile)
		},
	}

	newCmd.Flags().StringP("profile", "p", "", "Specify the profile you want to use")

	return newCmd
}

func setupNewProject(cmd *cobra.Command, args []string) (name, location string, profile *config.Profile, err error) {
	// Verify and set project location and name
	pname, location, err := setupProjectDirectory(args)
	if err != nil {
		return "", "", nil, err
	}

	setProjectNameLocation(pname, location)

	// Verify and set the user profile to use
	if cmd.Flag("profile").Value.String() != "" {
		viper.Set("profile", cmd.Flag("profile").Value.String())
	}

	profile, err = setProfile(cmd.Flag("profile").Value.String())

	return pname, location, profile, err
}

func setProfile(profileName string) (*config.Profile, error) {
	// If profile is given
	if profileName == "" {
		// No profile was specified, we're therefore calling the default profile
		fmt.Println("Loading default profile")

		profileName = viper.GetString(config.DefaultConfigProfileKeyName)

		if profileName == "" {
			return nil, errors.New("error : no profile was specified, and no default profile was found")
		}
	}

	profile, err := config.LoadProfile(profileName)
	if err == nil {
		viper.Set("profile", profileName)
	}

	return profile, err
}

func setupProjectDirectory(args []string) (pname, location string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", "", errors.Wrap(err, "unable to get working directory")
	}

	// If no argument was given, we develop the project inside the current directory, thus inheriting its name
	if len(args) == 0 {
		pname = path.Base(wd)
		location = wd
	} else {
		pname, location, err = newProjectDirectory(wd, args[0])
	}

	return pname, location, err
}

func newProjectDirectory(wd, argument string) (pname, location string, err error) {
	// Create Project destination folder
	if err := os.MkdirAll(argument, config.DirMode); err != nil {
		return "", "", errors.Wrapf(err, "Could not build project directory in '%s'", argument)
	}

	fmt.Printf("Build project directory %s\n", argument)

	if err := os.Chdir(argument); err != nil {
		return "", "", errors.Wrapf(err, "Could not change into project directory '%s'", argument)
	}

	return path.Base(argument), path.Join(wd, argument), nil
}

func setProjectNameLocation(pname, location string) {
	viper.Set("name", pname)
	viper.Set("location", location)
}

func newProject(name, location string, profile *config.Profile) {
	// Arguments are set and validated, initiated new project
	project := templates.NewProject(profile, name, location)

	// Build project
	if err := project.Build(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Project %s was successfully created.\n", name)
	os.Exit(0)
}
