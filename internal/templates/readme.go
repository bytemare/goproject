// Package templates holds the template and project building functions
package templates

const readmeIdentifier = "readme"

// readme implements the fileTemplate interface for a Readme
type readme struct {
	*file
	ProjectName string
	CI          *interface{}
}

// readmeConstructor returns the file content populated with the relevant values
func readmeConstructor(project *Project) (*file, error) {
	return newProjectFile(newReadme(project.Name))
}

func newReadme(projectName string) *readme {
	f, d, t := readmeValues()

	return &readme{
		file:        newFile(readmeIdentifier, f, d, t),
		ProjectName: projectName,
		CI:          nil,
	}
}

func (r *readme) getIdentifier() string {
	return r.identifier
}

func (r *readme) getFilename() string {
	return r.filename
}

func (r *readme) getTemplate() string {
	return r.template
}

func readmeValues() (f, d, t string) { //nolint:funlen // length is due to a constant, no complexity here
	const filename = "README"

	const directory = "."

	const template = `
[Template inspired by Jesse Luoto // https://github.com/jehna/readme-best-practices/blob/master/README-default.md]

# {{.ProjectName}}
![Logo of the Project](https://raw.githubusercontent.com/jehna/readme-best-practices/master/sample-logo.png)

  {{if .CI -}}
      {{- range .CI}}
  {{. -}}
      {{end}}
  {{end}}
> Additional information or tagline

A brief description of your Project, what it is used for and how does life get
awesome when someone starts to use it.

## Installing / Getting started

A quick introduction of the minimal setup you need to get a hello world up &
running.

\'\'\'shell
	packagemanager install {{.ProjectName}}
	{{.ProjectName}} start
	{{.ProjectName}} "Do something!"  # prints "Nah."
	\'\'\'

Here you should say what actually happens when you execute the code above.

### Initial Configuration

Some projects require initial configuration (e.g. access tokens or keys, \'npm i\').
This is the section where you would document those requirements.

## Developing

Here's a brief intro about what a developer must do in order to start developing
the Project further:

\'\'\'shell
	git clone https://github.com/your/{{.ProjectName}}.git
	cd {{.ProjectName}}/
		packagemanager install
	\'\'\'

And state what happens step-by-step.

### Building

If your Project needs some additional steps for the developer to build the
Project after some code changes, state them here:

\'\'\'shell
	./configure
	make
	make install
	\'\'\'

Here again you should state what actually happens when the code above gets
executed.

### Deploying / Publishing

In case there's some step you have to take that publishes this Project to a
server, this is the right time to state it.

\'\'\'shell
	packagemanager deploy {{.ProjectName}} -s server.com -u username -p password
	\'\'\'

And again you'd need to tell what the previous code actually does.

## Features

What's all the bells and whistles this Project can perform?
* What's the main functionality
* You can also do another thing
* If you get really randy, you can even do this

## Configuration

Here you should write what are all of the configurations a user can enter when
using the Project.

#### Argument 1
Type: \'String\'
Default: \''default value'\'

	State what an argument does and how you can use it. If needed, you can provide
	an example below.

		Example:
	\'\'\'bash
{{.ProjectName}} "Some other value"  # Prints "You're nailing this readme!"
\'\'\'

	#### Argument 2
Type: \'Number|Boolean\'
Default: 100

	Copy-paste as many of these as you need.

	## Usage

	> Some usage examples, with arguments and expected result

	## Supported Go versions

	> Go versions the Project has been tested and validated on

	Go Project works with the last three major Go versions, which are 1.11, 1.12 and 1.13 at the moment.

	## Contributing

	When you publish something open source, one of the greatest motivations is that
	anyone can just jump in and start contributing to your Project.

		These paragraphs are meant to welcome those kind souls to feel that they are
	needed. You should state something like:

	"If you'd like to contribute, please fork the repository and use a feature
	branch. Pull requests are warmly welcome."

	If there's anything else the developer needs to know (e.g. the code style
	guide), you should link it here. If there's a lot of things to take into
	consideration, it is common to separate this section to its own file called
	\'CONTRIBUTING.md\' (or similar). If so, you should say that it exists here.

	## Links

	Even though this information can be found inside the Project on machine-readable
	format like in a .json file, it's good to include a summary of most useful
	links to humans using your Project. You can include links like:

	- Project homepage: https://your.github.com/{{.ProjectName}}/
	- Repository: https://github.com/your/{{.ProjectName}}/
	- Issue tracker: https://github.com/your/{{.ProjectName}}/issues
	- In case of sensitive bugs like security vulnerabilities, please contact
my@email.com directly instead of using issue tracker. We value your effort
to improve the security and privacy of this Project!
- Related projects:
- Your other Project: https://github.com/your/other-Project/
- Someone else's Project: https://github.com/someones/{{.ProjectName}}/

## Licensing

One really important part: Give your Project a proper license. Here you should
state what the license is and how to find the text version of the license.
Something like:

"The code in this Project is licensed under MIT license."
`

	return filename, directory, template
}
