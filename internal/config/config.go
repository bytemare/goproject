package config

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/spf13/viper"
)

const (
	ConfDirName                 string = "goproject"
	ConfFileName                string = "conf.toml"
	ConfType                    string = "toml"
	ConfTitle                   string = "GoProject Configuration File"
	ProfileDirConfigKeyName     string = "profiledir"
	DefaultProfileConfigKeyName string = "defaultprofile"
	DefaultAutoUpdateKeyName    string = "autoupdate"
	DefaultAutoUpdateValue      string = "true"
)

func configDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(dir, ConfDirName), nil
}

//
func initConfig(location string) error {
	viper.SetConfigFile(path.Join(location, ConfFileName))
	viper.SetConfigType(ConfType)
	viper.SetConfigPermissions(0700)
	viper.Set("title", ConfTitle)
	viper.Set(ProfileDirConfigKeyName, ProfileDirName)
	viper.Set(DefaultProfileConfigKeyName, DefaultProfileName)
	viper.Set(DefaultAutoUpdateKeyName, true)
	if err := viper.WriteConfig(); err != nil {
		return err
	}

	// Create profile directory and default profile
	return InitProfiles(location)
}

func loadConfig() error {
	dir, _ := configDir()
	if dir != "" {
		viper.AddConfigPath(dir)
	}
	viper.SetConfigFile(path.Join(dir, ConfFileName))
	viper.SetConfigType(ConfType)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return errors.Wrap(err, "Config file could not be found")
		}
		return errors.Wrapf(err, "Config file exists but could not be loaded")
	}

	return nil
}

// getProfileDirName returns the path to the goproject profile directory
func GetProfileDirName() string {
	// Get configuration directory
	dir, _ := configDir()
	prof := fmt.Sprintf("%v", viper.Get(ProfileDirConfigKeyName))
	return path.Join(dir, prof)
}

func GetDefaultProfile() string {
	return fmt.Sprintf("%v", viper.Get(DefaultProfileConfigKeyName))
}

func Initialise() {
	// read in environment variables that match
	viper.AutomaticEnv()

	// Get configuration directory
	dir, err := configDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// If configuration directory does not exists, create and load it
	if _, err := os.Stat(dir); err != nil {
		// Configuration dir and file do not exist, so create them
		if err := os.MkdirAll(dir, 0700); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Create default config file
		if err := initConfig(dir); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return
	}

	// Check if config file exists
	configFile := path.Join(dir, ConfFileName)
	if _, err := os.Stat(configFile); err == nil {
		if err := loadConfig(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		// File does not exist
		if err := initConfig(dir); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
