package internal

import "github.com/spf13/viper"

type sonar struct {
	SonarOrg     string
	SonarProject string
	RepoCI       string
	RepoURL      string
}

func newSonar(sonarOrg, sonarProject, repoCI, repoURL string) *sonar {
	return &sonar{
		SonarOrg:     sonarOrg,
		SonarProject: sonarProject,
		RepoCI:       repoCI,
		RepoURL:      repoURL,
	}
}

func NewSonar(conf *viper.Viper) (*ProjectFile, error) {
	sonarOrg := conf.GetString("sonar.org")
	sonarProject := conf.GetString("name")
	repoCI := conf.GetString("travis.profile")
	repoURL := conf.GetString("git.URL")
	return NewProjectFile("sonar", "sonar-project.properties", newSonar(sonarOrg, sonarProject, repoCI, repoURL))
}

func (sonar) getTemplate() string {
	return sonarTemplate
}

const sonarTemplate = `# Project identification
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
# Path is relative to the sonar-project.properties file. Replace "\" by "/" on Windows.
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
