// Package templates holds the template and project building functions
package templates

const versionIdentifier = "version"

type version struct {
	*file
	PackageName string
}

// versionConstructor returns the file content populated with the relevant values
func versionConstructor(project *Project) (*file, error) { //nolint:unparam // project is not needed when no variables
	return newProjectFile(newVersion(project.Name))
}

func newVersion(packageName string) *version {
	f, d, t := versionValues()

	return &version{
		file:        newFile(versionIdentifier, f, d, t),
		PackageName: packageName,
	}
}

func (v *version) getIdentifier() string {
	return v.identifier
}

func (v *version) getFilename() string {
	return v.filename
}

func (v *version) getTemplate() string {
	return v.template
}

func versionValues() (f, d, t string) {
	const filename = "version.go"

	const directory = "internal/version"

	const template = `// Package version holds all the {{.PackageName}}'s core mechanisms
package version

import "fmt"

var (
	version = "?" //nolint:gochecknoglobals // package scoped / set at compile time to inject app version
	commit  = "?" //nolint:gochecknoglobals // package scoped / set at compile time to inject commit hash
	date    = "?" //nolint:gochecknoglobals // package scoped / set at compile time to inject compilation time
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
