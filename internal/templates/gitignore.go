// Package templates holds the template and project building functions
package templates

const gitignoreIdentifier = "gitignore"

// gitignoreConstructor returns the file content populated with the relevant values
func gitignoreConstructor(project *Project) (*file, error) { //nolint:unparam // project is not needed when no variables
	f, d, t := gitignoreValues()
	return newProjectFile(newFile(gitignoreIdentifier, f, d, t))
}

func gitignoreValues() (f, d, t string) {
	const filename = ".gitignore"

	const directory = "."

	const template = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, build with go test -c
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# IDE settings
.idea

# Directories
bin
coverage
`

	return filename, directory, template
}
