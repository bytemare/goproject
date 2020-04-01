// Package templates holds the template and project building functions
package templates

const docIdentifier = "doc"

type doc struct {
	*file
	PackageName string
}

// docConstructor returns the file content populated with the relevant values
func docConstructor(project *Project) (*file, error) {
	return newProjectFile(newDoc(project.Name))
}

func newDoc(packageName string) *doc {
	f, d, t := docValues()

	return &doc{
		file:        newFile(docIdentifier, f, d, t),
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

func docValues() (f, d, t string) {
	const filename = "doc.go"

	const directory = "."

	const template = `/*
Package {{.PackageName}} [short description]

*/
package {{.PackageName}}
`

	return filename, directory, template
}
