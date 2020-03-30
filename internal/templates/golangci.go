// Package templates holds the template and project building functions
package templates

// golangciConstructor returns the file content populated with the relevant values
func golangciConstructor(project *Project) (*file, error) { //nolint:unparam // project is not needed when no variables
	return newProjectFile(newFile(golangciIdentifier, golangciFilename, golangciTemplate))
}

const (
	golangciIdentifier = "golangci"
	golangciFilename   = ".golangci.yml"
	golangciTemplate   = `linters-settings:
  golint:
    min-confidence: 0
  maligned:
    suggest-new: true

linters:
  enable-all: true
  disable:
    - wsl
`
)
