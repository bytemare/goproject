package internal

type doc struct {
	PackageName string
}

func newDoc(packageName string) *doc {
	return &doc{PackageName: packageName}
}

func NewDoc(packageName string) (*ProjectFile, error) {
	return NewProjectFile("doc", "doc.go", newDoc(packageName))
}

func (doc) getTemplate() string {
	return `/*
Package {{.PackageName}} [short description]

*/
package {{.PackageName}}
`
}
