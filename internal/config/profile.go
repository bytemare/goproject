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

const (
	ProfileDirName string      = "profiles"
	DirMode        os.FileMode = 0700
	FileMode       os.FileMode = 0600
)

type Profile struct {
	//id     string
	Author *Author
	//layout map[string]map[string][]string
	//vcs    vcs
	//ci     ci
	Conf *viper.Viper
}

type Author struct {
	Name    string
	Contact string
}

/*type ci struct {
	platforms	[]string
	URL			map[string]string // indexes the project URL to the platform
}*/

/*type vcs struct {
	name	string // git
	repo	string // remote url to repo
}*/

// InitProfiles creates the profile directory if it does not exist
func InitProfiles(location string) error {
	profileDir := path.Join(location, ProfileDirName)
	if err := os.MkdirAll(profileDir, DirMode); err != nil {
		return err
	}

	// Write default profile
	return ioutil.WriteFile(path.Join(profileDir, DefaultProfileName), []byte(DefaultProfileContent), FileMode)
}

// loadProfile loads a profile from disk and returns a populated profile structure
func LoadProfile(name string) (*Profile, error) {
	fileProfile := viper.New()
	fileProfile.SetConfigFile(path.Join(GetProfileDirName(), name))
	fileProfile.SetConfigType(ConfType)
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
	dir := GetProfileDirName()

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(files))

	for _, f := range files {
		if !f.IsDir() {
			res = append(res, strings.TrimSuffix(f.Name(), "."+ConfType))
		}
	}

	return res, nil
}

// verifyProfile checks whether the profile is registered. If yes, returns the path to the profile configuration file
func verifyProfile(profile string) (string, error) {
	profiles, err := ListProfiles()
	if err != nil {
		return "", err
	}

	dir := GetProfileDirName()

	for _, p := range profiles {
		if p == profile {
			return path.Join(dir, p+"."+ConfType), nil
		}
	}
	return "", fmt.Errorf("error : profile '%s' was not found", profile)
}

// SetDefaultProfile sets the given profile as default
func SetDefaultProfile(profile string) error {
	if _, err := verifyProfile(profile); err != nil {
		return err
	}

	var profileFile string
	if strings.HasSuffix(profile, "."+ConfType) {
		profileFile = profile
	} else {
		profileFile = profile + "." + ConfType
	}

	viper.Set(DefaultProfileConfigKeyName, profileFile)
	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}

func CreateNewProfileFile(profile string) (string, error) {
	profilePath, _ := verifyProfile(profile)
	if profilePath != "" {
		return "", fmt.Errorf("profile %s already exists", profile)
	}

	dir := GetProfileDirName()
	if strings.HasSuffix(profile, "."+ConfType) {
		profilePath = path.Join(dir, profile)
	} else {
		profilePath = path.Join(dir, profile+"."+ConfType)
	}

	if err := ioutil.WriteFile(profilePath, []byte(NewProfileContent), FileMode); err != nil {
		return "", err
	}

	return profilePath, nil
}

// nanoProfile opens a text editor to the specified profile configuration file
func NanoProfile(profile string) error {
	profilePath, err := verifyProfile(profile)
	if err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		return errors.New("windows is not yet supported")
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("nano", profilePath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err = cmd.Start(); err != nil {
			return errors.Wrap(err, "nano could not be started")
		}

		if err = cmd.Wait(); err != nil {
			return errors.Wrap(err, "error while editing profile file")
		}

		return nil
	}

	return errors.New("you're running on an unsupported system")
}
