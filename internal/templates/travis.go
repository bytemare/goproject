// Package templates holds the template and project building functions
package templates

const travisIdentifier = "travis"

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
	f, d, t := travisValues()

	return &travis{
		file:     newFile(travisIdentifier, f, d, t),
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

func travisValues() (f, d, t string) {
	const filename = ".travis.yml"

	const directory = "."

	const template = `language: go

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
  - "Unit Tests and Coverage"
  #- name: deploy
  #  if: branch = release

addons:
  apt:
    packages:
      - "python3"
      - "python3-dev"
      - "python3-pip"
      - "python3-setuptools"

jobs:
  include:
    - stage: "Static Analysis, Unit Tests and Coverage"
      go: 1.14.x
      name: "GolangCI Linting, and Snyk Analysis"
      os: linux
      install:
        - make prepare-pre-commit
        - npm install -g snyk
      script:
        - pre-commit autoupdate
        - pre-commit run --all-files
        - snyk test
      after_success:
        - snyk monitor
    - go: 1.14.x
      name: "Unit Tests and Coverage"
      {{if .SonarOrg -}}
      addons:
        sonarcloud:
          organization: "{{.SonarOrg}}"
          token:
            secure: ${SONAR_TOKEN}
	  {{end}}
      os: linux
      install:
        - make prepare-tests
      script:
        - make cover
      after_success:
        - sonar-scanner -X
        #- bash <(curl -s https://codecov.io/bash)
        #- goveralls -coverprofile=coverage.out -service=travis-ci
    #- stage: release
    #  name: "Release a new version"
    #  deploy:
    #    provider: script
    #    cleanup: true
    #    script:
    #      - nvm install lts/*
    #      - npx semantic-release
    #    on:
    #      all_branches: true

go:
  - 1.11.x
  - 1.12.x
  - 1.13.x
  - 1.14.x
os:
  - linux
  - osx
  - windows
script:
  - if [ "$TRAVIS_OS_NAME" = "windows" ];
    then go test -v -i -race -covermode=atomic;
    else make test;
    fi

`

	return filename, directory, template
}
