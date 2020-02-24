package internal

import (
	"github.com/spf13/viper"
)

type travis struct {
	RepoURL  string
	SonarOrg string
}

//func newTravis(repoURL, sonarOrg string) *travis {
func newTravis(repoURL, sonarOrg string) *travis {
	return &travis{
		RepoURL:  repoURL,
		SonarOrg: sonarOrg,
	}
}

func NewTravis(conf *viper.Viper) (*ProjectFile, error) {
	sonarOrg := conf.GetString("sonar.org")
	repoURL := conf.GetString("git.URL")
	return NewProjectFile("travis", ".travis.yml", newTravis(repoURL, sonarOrg))
}

func (travis) getTemplate() string {
	return travisTemplate
}

const travisTemplate = `language: go

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
