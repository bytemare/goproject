package internal

import (
	"bytes"
	"io/ioutil"
	"path"
	"text/template"

	"github.com/pkg/errors"
)

// ProjectFile represents a project's file
type ProjectFile struct {
	name     string
	filename string
	location string
	content  string
}

// templateValues are structures holding the values for the corresponding template,
// with a function to get the raw template
type templateValues interface {
	getTemplate() string // returns the raw template
}

func newProjectFile(name, file, content string) *ProjectFile {
	return &ProjectFile{
		name:     name,
		filename: file,
		content:  content,
	}
}

//func NewProjectFile(name, filename string, fn newTemplateValuesFunc, args ...string) (*ProjectFile, error) {
func NewProjectFile(name, filename string, values templateValues) (*ProjectFile, error) {
	// Build a project file
	pf := newProjectFile(name, filename, "")

	// Process the template with values
	if err := pf.process(values); err != nil {
		return nil, err
	}

	return pf, nil
}

func (pf *ProjectFile) process(vals templateValues) error {
	// Parse the template
	tmpl, err := template.New(pf.name).Parse(vals.getTemplate())
	if err != nil {
		return errors.Wrapf(err, "could not parse template string from '%s'", pf.name)
	}

	// Process the template with the associated values
	var buff bytes.Buffer
	if err := tmpl.Execute(&buff, vals); err != nil {
		return errors.Wrapf(err, "could not execute template for '%s' with values of '%v'", pf.name, vals)
	}

	pf.content = buff.String()
	return nil
}

func (pf *ProjectFile) Write() error {
	return ioutil.WriteFile(path.Join(pf.location, pf.name), []byte(pf.content), 0600)
}
