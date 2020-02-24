package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bytemare/goproject/internal/config"
	templates "github.com/bytemare/goproject/internal/templates"
	"github.com/pkg/errors"
)

type Project struct {
	Profile *config.Profile
	Name    string
	Path    string
	Layout  Layout
	Author  *config.Author
}

type Layout struct {
	Directories []string
	Files       []string
}

func NewProject(prof *config.Profile, name, location string) *Project {
	project := &Project{
		Profile: prof,
		Name:    name,
		Path:    filepath.Join(location, name),
		Layout: Layout{
			Directories: make([]string, 0),
			Files:       make([]string, 0),
		},
		Author: prof.Author,
	}

	err := prof.Conf.Unmarshal(project)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	}
	return project
}

// Build creates the project structure and writes files
func (p *Project) Build() error {
	// Create project destination folder if it does not exist already
	if err := os.MkdirAll(p.Path, 0600); err != nil {
		return err
	}

	if err := os.Chdir(p.Path); err != nil {
		return errors.Wrapf(err, "could not build project in '%s'", p.Path)
	}

	// Build the directory and file layout
	p.buildDirs()
	p.buildFiles()

	// Initialise git and go modules
	if err := goMod(); err != nil {
		return err
	}
	if err := gitInit(); err != nil {
		return err
	}

	return nil
}

func (p *Project) buildDirs() {
	if len(p.Layout.Directories) == 0 {
		return
	}

	fmt.Println("Creating directory layout.")
	for i, d := range p.Layout.Directories {
		fmt.Printf("\t> %d : Building directory %s ... ", i, d)
		if err := os.MkdirAll(p.Path, 0700); err != nil {
			fmt.Printf("error : %v\n", err)
		} else {
			fmt.Printf("success.\n")
		}
	}
}

func (p *Project) buildFiles() {
	if len(p.Layout.Files) == 0 {
		return
	}

	fmt.Println("Creating files.")
	for i, f := range p.Layout.Files {
		fmt.Printf("\t> %d : Build file %s ... ", i, f)
		if pf, err := p.buildFile(f); err != nil {
			fmt.Printf("error : %v\n", err)
		} else {
			if err := pf.Write(); err != nil {
				fmt.Printf("error : %v\n", err)
			}
			fmt.Printf("success.\n")
		}
	}
}

// buildFile() writes the template
func (p *Project) buildFile(fileID string) (pf *templates.ProjectFile, err error) {
	switch fileID {
	case "doc":
		pf, err = templates.NewDoc(p.Name)
	case "readme":
		pf, err = templates.NewReadme(p.Name)
	case "makefile":
		pf, err = templates.NewMakefile()
	case "dockerfile":
		pf, err = templates.NewDockerfile(p.Author.Name+" "+p.Author.Contact, p.Name)
	case "gitignore":
		pf, err = templates.NewGitIgnore()
	case "golangci":
		pf, err = templates.NewGolangCI()
	case "sonar":
		pf, err = templates.NewSonar(p.Profile.Conf)
	case "travis":
		pf, err = templates.NewTravis(p.Profile.Conf)
	default:
		err = fmt.Errorf("error : '%s' is not a registered project file", fileID)
	}

	if err != nil {
		return nil, err
	}

	return pf, nil
}

func gitInit() error {
	fmt.Println("Initialising git.")
	return runCmd("git", "init")
}

func goMod() error {
	fmt.Println("Initialising go modules.")
	return runCmd("go", "mod", "init")
}

func runCmd(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	return cmd.Run()
}
