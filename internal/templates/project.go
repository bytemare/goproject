// Package templates holds the template and project building functions
package templates

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/bytemare/goproject/internal/config"
)

// Project structure contains all the necessary information regarding a Project,
// and a reference to the profile it is build for
type Project struct {
	Profile *config.Profile
	Name    string
	Path    string
	Layout  layout
	Author  *config.Author
}

// layout represents the directory and file layout of the Project
type layout struct {
	Directories []string
	Files       []string
}

// NewProject returns a new Project structure given a name, where it is to be created,
// and a profile containing the directives for the Project layout
func NewProject(prof *config.Profile, name, location string) *Project {
	project := &Project{
		Profile: prof,
		Name:    name,
		Path:    location,
		Layout: layout{
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

// Build creates the Project structure and writes files
func (p *Project) Build() error {
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

const buildErrFormat = "error : %v\n"

func (p *Project) buildDirs() {
	if len(p.Layout.Directories) == 0 {
		return
	}

	fmt.Println("Creating directory layout.")

	for i, d := range p.Layout.Directories {
		fmt.Printf("\t> %d : Building directory %s ... ", i, d)

		if err := os.MkdirAll(d, config.DirMode); err != nil {
			fmt.Printf(buildErrFormat, err)
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

	for i, fid := range p.Layout.Files {
		fmt.Printf("\t> %d : Build file %s ... ", i, fid)

		if f, err := p.buildFile(fid); err != nil {
			fmt.Printf(buildErrFormat, err)
		} else {
			e, err := exists(path.Join(f.directory, f.filename))
			if e {
				if err == nil {
					fmt.Println("already exists. Skipping.")
				}
				continue
			}

			if err := f.write(); err != nil {
				fmt.Printf(buildErrFormat, err)
			}
			fmt.Printf("success.\n")
		}
	}
}

func (p *Project) buildFile(fileID string) (*file, error) {
	// fetch the constructor for corresponding file identifier
	constructor, err := getFileConstructor(fileID)
	if err != nil {
		return nil, err
	}

	// execute the constructor and return its output
	return constructor(p)
}

func gitInit() error {
	fmt.Println("Initialising git.")

	e, err := exists(".git")
	if e {
		if err == nil {
			fmt.Println("\tGit directory (.git) already exists. Skipping initialisation.")
		}

		return err
	}

	cmd := exec.Command("git", "init")

	return cmd.Run()
}

func goMod() error {
	fmt.Println("Initialising go modules.")

	cmd := exec.Command("go", "mod", "init")

	return cmd.Run()
}

// exists returns whether the given file or directory exists
func exists(location string) (bool, error) {
	_, err := os.Stat(location)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}
