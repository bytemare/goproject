// Package templates holds the template and project building functions
package templates

type doc struct {
	*file
	PackageName string
}

// docConstructor returns the file content populated with the relevant values
func docConstructor(project *Project) (*file, error) {
	return newProjectFile(newDoc(project.Name))
}

func newDoc(packageName string) *doc {
	return &doc{
		file:        newFile(docIdentifier, docFilename, docTemplate),
		PackageName: packageName,
	}
}

func (d *doc) getIdentifier() string {
	return d.identifier
}

func (d *doc) getFilename() string {
	return d.filename
}

func (d *doc) getTemplate() string {
	return d.template
}

const (
	docIdentifier = "doc"
	docFilename   = "doc.go"
	docTemplate   = `/*
Package {{.PackageName}} [short description]

*/
package {{.PackageName}}
`
)
