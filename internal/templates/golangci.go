package internal

type golangci struct {
}

func newGolangCI() *golangci {
	return &golangci{}
}

func NewGolangCI() (*ProjectFile, error) {
	return NewProjectFile("golangci", ".golangci.yml", newGolangCI())
}

func (golangci) getTemplate() string {
	return `linters-settings:
  golint:
    min-confidence: 0
  maligned:
    suggest-new: true

linters:
  enable-all: true
  disable:
    - wsl
`
}
