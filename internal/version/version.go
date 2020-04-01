// Package version holds all the goproject's core mechanisms
package version

import "fmt"

var (
	version = "?" //nolint:gochecknoglobals // used at compile time to inject app version
	commit  = "?" //nolint:gochecknoglobals // used at compile time to inject commit hash
	date    = "?" //nolint:gochecknoglobals // used at compile time to inject compilation time
)

// GetVersion returns the version the program was compiled with
func GetVersion() string {
	return version
}

// GetCommit returns the commit the program was compiled with
func GetCommit() string {
	return commit
}

// PrintableVersion returns a string representing the version, commit tag and date of built
func PrintableVersion() string {
	return fmt.Sprintf("Version %s:%s - compiled on %s", version, commit, date)
}

// Upgrade attempts to upgrade the program to the latest version
func Upgrade() {
	fmt.Printf("Upgrade not implemented yet - %s\n", PrintableVersion())
}
