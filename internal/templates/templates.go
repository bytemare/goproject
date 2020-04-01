// Package templates holds the template and project building functions
package templates

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"text/template"

	"github.com/pkg/errors"

	"github.com/bytemare/goproject/internal/config"
)

// fileTemplate are structures holding basic values for the file and corresponding template
type fileTemplate interface {
	getIdentifier() string // returns the identifier for a template, e.g. "dockerfile"
	getFilename() string   // returns the filename
	getDirectory() string  // returns the directory of the file
	getTemplate() string   // returns the raw template
	getFile() *file        // returns the potentially embedded file struct
}

type file struct {
	identifier string
	filename   string
	directory  string
	template   string
	content    string
}

func (f *file) getIdentifier() string {
	return f.identifier
}

func (f *file) getFilename() string {
	return f.filename
}

func (f *file) getDirectory() string {
	return f.directory
}

func (f *file) getTemplate() string {
	return f.template
}

func (f *file) getFile() *file {
	return f
}

func (f *file) write() error {
	return ioutil.WriteFile(path.Join(f.directory, f.filename), []byte(f.content), config.FileMode)
}

func newFile(identifier, filename, directory, rawTemplate string) *file {
	return &file{
		identifier: identifier,
		filename:   filename,
		directory:  directory,
		template:   rawTemplate,
	}
}

// getFileConstructor returns a registered projectFile constructor associated with the given file identifier
func getFileConstructor(fileID string) (constructor func(*Project) (*file, error), err error) {
	// Reference all available templates here
	registeredTemplates := map[string]func(*Project) (*file, error){
		docIdentifier:        docConstructor,
		dockerfileIdentifier: dockerfileConstructor,
		gitignoreIdentifier:  gitignoreConstructor,
		golangciIdentifier:   golangciConstructor,
		makefileIdentifier:   makefileConstructor,
		precommitIdentifier:  precommitConstructor,
		readmeIdentifier:     readmeConstructor,
		sonarIdentifier:      sonarConstructor,
		travisIdentifier:     travisConstructor,
		versionIdentifier:    versionConstructor,
	}

	constructor, ok := registeredTemplates[fileID]
	if !ok {
		err = fmt.Errorf("error : '%s' is not a registered Project file", fileID)
	}

	return constructor, err
}

func newProjectFile(ft fileTemplate) (*file, error) {
	// Parse the template
	tmpl, err := template.New(ft.getIdentifier()).Parse(ft.getTemplate())
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse template string from '%s'", ft.getIdentifier())
	}

	// Process the template with the associated values
	var buff bytes.Buffer
	if err := tmpl.Execute(&buff, ft); err != nil {
		return nil, errors.Wrapf(err, "could not execute template for '%s' with values of '%v'", ft.getIdentifier(), ft)
	}

	f := ft.getFile()
	f.content = buff.String()

	return f, nil
}
