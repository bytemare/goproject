// Package config groups the configuration and profile management mechanism
package config

import "os"

const (
	confDirName             string = "goproject"
	confFileName            string = "conf.toml"
	confType                string = "toml"
	confTitle               string = "GoProject Configuration File"
	profileDirConfigKeyName string = "profiledir"
	profileDirName          string = "profiles"
)

const (
	// DefaultTargetProjectLocation is the default location to deploy the new project in
	DefaultTargetProjectLocation string = "."

	// DefaultConfigProfileKeyName is the application's configuration entry label for the default profile to use
	DefaultConfigProfileKeyName string = "defaultprofile"

	// DefaultAutoUpdateKeyName is the application's configuration entry label for auto-updates
	DefaultAutoUpdateKeyName string = "autoupdate"
	// DefaultAutoUpdateValue is the default value for auto-updates
	DefaultAutoUpdateValue string = "true"

	// DirMode is the default file permissions mode
	DirMode os.FileMode = 0700
	// FileMode is the default directory permissions mode
	FileMode os.FileMode = 0600
)
