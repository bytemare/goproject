// Package templates holds the template and project building functions
package templates

type travis struct {
	*file
	RepoURL  string
	SonarOrg string
}

// travisConstructor returns the file content populated with the relevant values
func travisConstructor(project *Project) (*file, error) {
	conf := project.Profile.Conf
	sonarOrg := conf.GetString("sonar.org")
	repoURL := conf.GetString("git.URL")

	return newProjectFile(newTravis(repoURL, sonarOrg))
}

func newTravis(repoURL, sonarOrg string) *travis {
	return &travis{
		file:     newFile(travisIdentifier, travisFilename, travisTemplate),
		RepoURL:  repoURL,
		SonarOrg: sonarOrg,
	}
}

func (t *travis) getIdentifier() string {
	return t.identifier
}

func (t *travis) getFilename() string {
	return t.filename
}

func (t *travis) getTemplate() string {
	return t.template
}

const (
	travisIdentifier = "travis"
	travisFilename   = ".travis.yml"
	travisTemplate   = `language: go

env:
  - GO111MODULE=on

#branches:
  #only:
    #- dev

git:
  depth: false # Sonar doesn't like shallow clones

notifications:
  email: false

stages:
  - "Static Analysis, Unit Tests and Coverage"
  - test
  #- name: deploy
  #    if: branch = master

jobs:
  include:
    - stage: "Static Analysis, Unit Tests and Coverage"
      go: 1.13.x
      name: "GolangCI Linting and Snyk Analysis"
      os: linux
      install:
        - go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
        - npm install -g snyk
      script:
        - golangci-lint run -v ./...
        - snyk test
      after_success:
        - snyk monitor
    - go: 1.13.x
      name: "Unit Tests and Coverage"
	  {{if .SonarOrg -}}
      addons:
        sonarcloud:
          organization: "{{.SonarOrg}}"
          token:
            secure: ${SONAR_TOKEN}
	  {{end}}
      os: linux
      script:
        - go test -v -race -coverprofile=coverage.out -covermode=atomic
      after_success:
        - sonar-scanner
        - bash <(curl -s https://codecov.io/bash)
        - goveralls -coverprofile=coverage.out -service=travis-ci

go:
  - 1.11.x
  - 1.12.x
  - 1.13.x
os:
  - linux
  - osx
  - windows
script:
  - go test -v -race
`
)
