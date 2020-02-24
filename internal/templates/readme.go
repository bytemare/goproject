package internal

// readme implements the templateValues interface for a Readme
type readme struct {
	ProjectName string
	CI          *interface{}
}

func newReadme(projectName string) *readme {
	return &readme{
		ProjectName: projectName,
		CI:          nil,
	}
}

func NewReadme(projectName string) (*ProjectFile, error) {
	return NewProjectFile("readme", "README.md", newReadme(projectName))
}

func (readme) getTemplate() string {
	return `# {{.ProjectName}}

{{if .CI -}}
    {{- range .CI}}
{{. -}}
    {{end}}
{{end}}

> Short description

## Installation

> Installation instructions

## Usage

> Some usage examples, with arguments and expected result

## Supported Go versions

> Go versions the project has been tested and validated on

Go project works with the last three major Go versions, which are 1.11, 1.12 and 1.13 at the moment.

## Licence
`
}
