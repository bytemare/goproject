// Package templates holds the template and project building functions
package templates

const precommitIdentifier = "pre-commit"

// makefileConstructor returns the file content populated with the relevant values
func precommitConstructor(project *Project) (*file, error) { //nolint:unparam // project is not needed when no variables
	f, d, t := precommitValues()
	return newProjectFile(newFile(precommitIdentifier, f, d, t))
}

func precommitValues() (f, d, t string) {
	const filename = ".pre-commit-config.yaml"

	const directory = "."

	const template = `repos:
  -   repo: https://github.com/pre-commit/pre-commit-hooks
      rev: v2.5.0
      hooks:
        -   id: trailing-whitespace
        -   id: check-docstring-first
        -   id: end-of-file-fixer
        -   id: check-yaml
        -   id: check-toml
        -   id: detect-aws-credentials
            args: [--allow-missing-credentials]
        -   id: detect-private-key
        -   id: no-commit-to-branch
        -   id: check-added-large-files
  - repo: https://github.com/syntaqx/git-hooks
    rev: v0.0.16
    hooks:
      - id: forbid-binary
      - id: go-mod-tidy
      - id: shfmt
  - repo: git://github.com/detailyang/pre-commit-shell
    rev: 1.0.5
    hooks:
      - id: shell-lint
        args: [--format=json]
  - repo: https://github.com/mattlqx/pre-commit-sign
    rev: v1.1.1
    hooks:
      - id: sign-commit
  - repo: https://github.com/golangci/golangci-lint
    rev: v1.24.0
    hooks:
      - id: golangci-lint
  - repo: local
    hooks:
      - id: testing
        name: 'Testing'
        entry: make test
        language: 'system'
        files: Makefile
        description: "Runs the Makefile's test command"

`

	return filename, directory, template
}
