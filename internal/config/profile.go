// Package config groups the configuration and profile management mechanism
package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/pkg/errors"

	"github.com/spf13/viper"
)

var _ = map[string]string{
	"travis":   "https://travis-ci.com",
	"circleci": "https://circle-ci.com",
	"appveyor": "https://appveyor.com",
	"gitlab":   "https://gitlab.com",
}

// Profile associates an Author to a configuration file, describing a desired project layout
type Profile struct {
	Author *Author
	Conf   *viper.Viper
}

// Author represents the developer and holds their contact information
type Author struct {
	Name    string
	Contact string
}

// initProfiles creates the profile directory if it does not exist
func initProfiles(location string) error {
	profileDir := path.Join(location, profileDirName)
	if err := os.MkdirAll(profileDir, DirMode); err != nil {
		return err
	}

	// Write default profile
	return ioutil.WriteFile(path.Join(profileDir, defaultProfileName), []byte(defaultProfileContent), FileMode)
}

// LoadProfile loads a profile from disk and returns a populated profile structure
func LoadProfile(name string) (*Profile, error) {
	dir, err := GetProfileDirName()
	if err != nil {
		return nil, errors.Wrap(err, "Could not load profile")
	}

	fileProfile := viper.New()
	fileProfile.SetConfigFile(path.Join(dir, name))
	fileProfile.SetConfigType(confType)

	if err := fileProfile.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, errors.Wrapf(err, "Profile '%s' does not exist (or profile file not found)", name)
		}

		return nil, errors.Wrapf(err, "File for profile '%s' exists but could not be loaded", name)
	}

	var p Profile
	if err := fileProfile.Unmarshal(&p); err != nil {
		return nil, errors.Wrapf(err, "Could not load profile '%s'", name)
	}

	p.Conf = fileProfile

	return &p, nil
}

// ListProfiles returns the list of available profiles
func ListProfiles() ([]string, error) {
	dir, err := GetProfileDirName()
	if err != nil {
		return nil, errors.Wrap(err, "Could not list profiles")
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "Could not list profiles")
	}

	res := make([]string, 0, len(files))

	for _, f := range files {
		if !f.IsDir() {
			res = append(res, strings.TrimSuffix(f.Name(), "."+confType))
		}
	}

	return res, nil
}

// verifyProfile checks whether the profile is registered. If yes, returns the path to the profile configuration file,
// if not, returns an empty string and no error
func verifyProfile(profile string) (string, error) {
	profiles, err := ListProfiles()
	if err != nil {
		return "", err
	}

	dir, err := configDir()
	if err != nil {
		return "", errors.Wrap(err, "Could not verify profile")
	}

	for _, p := range profiles {
		if p == profile {
			return path.Join(dir, profileDirName, p+"."+confType), nil
		}
	}

	// profile was not found
	return "", nil
}

// SetDefaultProfile sets the given profile as default
func SetDefaultProfile(profile string) error {
	if _, err := verifyProfile(profile); err != nil {
		return err
	}

	var profileFile string

	if strings.HasSuffix(profile, "."+confType) {
		profileFile = profile
	} else {
		profileFile = profile + "." + confType
	}

	viper.Set(DefaultConfigProfileKeyName, profileFile)

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}

// CreateNewProfileFile creates a new file for the given profile name if it does not exist
func CreateNewProfileFile(profile string) (string, error) {
	p, err := verifyProfile(profile)
	if p != "" {
		return "", fmt.Errorf("profile %s already exists", profile)
	}

	if err != nil {
		return "", errors.Wrap(err, " could not create new profile")
	}

	dir, err := configDir()
	if err != nil {
		return "", errors.Wrap(err, "Could not create new profile file")
	}

	profilePath := ""

	if strings.HasSuffix(profile, "."+confType) {
		profilePath = path.Join(dir, profileDirName, profile)
	} else {
		profilePath = path.Join(dir, profileDirName, profile+"."+confType)
	}

	if err := ioutil.WriteFile(profilePath, []byte(newProfileContent), FileMode); err != nil {
		return "", err
	}

	return profilePath, nil
}

// EditProfile opens a text editor to the specified profile configuration file
func EditProfile(profile string) error {
	profilePath, err := verifyProfile(profile)
	if err != nil {
		return err
	}

	// security : verify the file indicated is not a symlink
	if _, err = isSymlink(profilePath); err != nil {
		return errors.Wrapf(err, "can't edit profile")
	}

	fi, err := os.Stat(profilePath)
	if err != nil {
		return errors.Wrapf(err, "can't edit profile")
	}

	if fi == nil {
		return fmt.Errorf("can't edit profile ' : os.Stat returned empty FileInfo()")
	}

	lastModified := fi.ModTime()

	// map editor and command line flags to OS
	editors := map[string]string{
		"linux":   "nano",
		"darwin":  "nano",
		"freebsd": "ee",
		"netbsd":  "vi",
		"openbsd": "ed",
		"plan9":   "sam",
	}

	v, ok := editors[runtime.GOOS]
	if !ok {
		return fmt.Errorf("%s is not yet supported for CLI profile editing", runtime.GOOS)
	}

	cmd := exec.Command(v, profilePath) //nolint:gosec // os/exec doesn't call shell, and args are limited
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Start(); err != nil {
		return errors.Wrapf(err, "%v could not be started", v[0])
	}

	if err = cmd.Wait(); err != nil {
		return errors.Wrap(err, "error while editing profile file")
	}

	fi, err = os.Stat(profilePath)
	if err != nil {
		fmt.Println(errors.Wrapf(err, "an error occurred while operating os.Stat() on file '%s'", profilePath))
	}

	if fi == nil {
		fmt.Println(fmt.Errorf("os.Stat returned empty FileInfo() on '%s'", profilePath))
	}

	modified := fi.ModTime()
	if modified == lastModified {
		fmt.Println("Warning: the profile file has not been modified.")
	}

	return nil
}

// isSymlink returns whether the file pointed to by path is a symlink. If yes, the error will contain more information
func isSymlink(filePath string) (bool, error) {
	fi, err := os.Lstat(filePath)
	if err != nil {
		return false, err
	}

	// Will return true if it is a symlink
	if fi.Mode()&os.ModeSymlink != 0 {
		link, err := os.Readlink(fi.Name())
		if err != nil {
			return true, err
		}

		return true, fmt.Errorf("%v is a symlink to %v", filePath, link)
	}

	return false, nil
}
