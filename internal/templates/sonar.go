// Package templates holds the template and project building functions
package templates

const sonarIdentifier = "sonar"

type sonar struct {
	*file
	SonarOrg     string
	SonarProject string
	RepoCI       string
	RepoURL      string
}

// sonarConstructor returns the file content populated with the relevant values
func sonarConstructor(project *Project) (*file, error) {
	conf := project.Profile.Conf
	sonarOrg := conf.GetString("sonar.org")
	sonarProject := conf.GetString("name")
	repoCI := conf.GetString("travis.profile")
	repoURL := conf.GetString("git.URL")

	return newProjectFile(newSonar(sonarOrg, sonarProject, repoCI, repoURL))
}

func newSonar(sonarOrg, sonarProject, repoCI, repoURL string) *sonar {
	f, d, t := sonarValues()

	return &sonar{
		file:         newFile(sonarIdentifier, f, d, t),
		SonarOrg:     sonarOrg,
		SonarProject: sonarProject,
		RepoCI:       repoCI,
		RepoURL:      repoURL,
	}
}

func (s *sonar) getIdentifier() string {
	return s.identifier
}

func (s *sonar) getFilename() string {
	return s.filename
}

func (s *sonar) getTemplate() string {
	return s.template
}

func sonarValues() (f, d, t string) {
	const filename = "sonar-Project.properties"

	const directory = "."

	const template = `# Project identification
sonar.organization={{.SonarOrg}}
sonar.projectKey={{.SonarProject}}
sonar.projectName={{.SonarProject}}
#sonar.projectVersion=1

# Project Metadata
sonar.links.ci={{.RepoCI}}
sonar.links.homepage={{.RepoURL}}
sonar.links.scm={{.RepoURL}}
sonar.host.url=https://sonarcloud.io

# Project files
# Path is relative to the sonar-Project.properties file. Replace "\" by "/" on Windows.
# This property is optional if sonar.modules is set.
sonar.sources=.
sonar.sourceEncoding=UTF-8
sonar.exclusions=**/*_test.go,**/vendor/**

# Testing
#sonar.tests=./Tests/
sonar.test.inclusions=**/*_test.go
sonar.go.tests.reportPaths=./trace.out
#sonar.test.exclusions=**/vendor/**

# Coverage
sonar.coverage.exclusions= **/*_test.go
sonar.go.coverage.reportPaths=./coverage.out

# Other Modules
#sonar.go.golangci-lint.reportPaths=
`

	return filename, directory, template
}
