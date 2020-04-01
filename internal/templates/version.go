// Package templates holds the template and project building functions
package templates

const versionIdentifier = "version"

// docConstructor returns the file content populated with the relevant values
func versionConstructor(project *Project) (*file, error) { //nolint:unparam // project is not needed when no variables
	f, d, t := versionValues()
	return newProjectFile(newFile(versionIdentifier, f, d, t))
}

func versionValues() (f, d, t string) {
	const filename = "version.go"

	const directory = "internal/version"

	const template = `// Package version holds all the goproject's core mechanisms
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

`

	return filename, directory, template
}
