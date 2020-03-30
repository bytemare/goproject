// Package config groups the goproject configuration and profile management mechanism
package config

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/spf13/viper"
)

func configDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return path.Join(dir, confDirName), nil
}

//
func initConfig(location string) error {
	fm := FileMode

	viper.SetConfigFile(path.Join(location, confFileName))
	viper.SetConfigType(confType)
	viper.SetConfigPermissions(fm)
	viper.Set("title", confTitle)
	viper.Set(profileDirConfigKeyName, profileDirName)
	viper.Set(DefaultConfigProfileKeyName, defaultProfileName)
	viper.Set(DefaultAutoUpdateKeyName, true)

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	// Create profile directory and default profile
	return initProfiles(location)
}

func loadConfig() error {
	dir, err := configDir()
	if err != nil {
		return errors.Wrap(err, "Could not load configuration")
	}

	viper.AddConfigPath(dir)
	viper.SetConfigFile(path.Join(dir, confFileName))
	viper.SetConfigType(confType)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return errors.Wrap(err, "Config file could not be found")
		}

		return errors.Wrapf(err, "Config file exists but could not be loaded")
	}

	return nil
}

// GetProfileDirName returns the path to the goproject profile directory
func GetProfileDirName() (string, error) {
	// Get configuration directory
	dir, err := configDir()
	if err != nil {
		return "", errors.Wrap(err, "Could not get profile directory")
	}

	prof := fmt.Sprintf("%v", viper.Get(profileDirConfigKeyName))

	return path.Join(dir, prof), nil
}

/*
func GetDefaultProfile() string {
	return fmt.Sprintf("%v", viper.Get(DefaultProfileConfigKeyName))
}*/

// Initialise reads or creates configurations in order to aliment the program with needed input
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
		if err := os.MkdirAll(dir, DirMode); err != nil {
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
	configFile := path.Join(dir, confFileName)
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
