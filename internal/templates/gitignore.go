package internal

type gitignore struct {
}

func newGitignore() *gitignore {
	return &gitignore{}
}

func NewGitIgnore() (*ProjectFile, error) {
	return NewProjectFile("gitignore", ".gitignore", newGitignore())
}

func (gitignore) getTemplate() string {
	return `# Binaries for programs and plugins
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
}
