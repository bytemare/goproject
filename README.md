# goproject
Set up new go projects with personalised layouts and configuration files for your CI environment.

[![Build Status](https://travis-ci.com/bytemare/goproject.svg?branch=master)](https://travis-ci.com/bytemare/goproject)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=goproject&metric=alert_status)](https://sonarcloud.io/dashboard?id=goproject)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=goproject&metric=coverage)](https://sonarcloud.io/dashboard?id=goproject)
[![Go Report Card](https://goreportcard.com/badge/github.com/bytemare/goproject)](https://goreportcard.com/report/github.com/bytemare/goproject)
[![GolangCI](https://golangci.com/badges/github.com/bytemare/goproject.svg)](https://golangci.com/r/github.com/bytemare/goproject)
[![GoDoc](https://godoc.org/github.com/bytemare/goproject?status.svg)](https://godoc.org/github.com/bytemare/goproject)
[![fuzzit](https://app.fuzzit.dev/badge?org_id=bytemare-gh)](https://app.fuzzit.dev/orgs/bytemare-gh/dashboard)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)


Concentrate on writing code for your project, and don't waste time setting it up, goproject does it for you.

Set up new go projects with personalised layouts and configuration files for your CI environment.
Available templates contain advanced CI configuration files, a Dockerfile, and some guidelines. Some example package functions and tests are also given.

- [GoProject](#goproject)
    - [Installation](#install)
        - [Binary](#binary)
        - [macOS](#macos)
        - [Docker](#docker)
        - [Go](#go)
    - [Configuration](#config)
    - [Usage](#use)
    - [Changelog](#changelog)
    - [Features](#features)
    - [Thanks](#thanks)
    - [Licence](#licence)

## Install

### Binary

To get the latest auto-updatable stable [release available](https://github.com/bytemare/goproject/releases).

```bash
# binary will be $(go env GOPATH)/bin/golangci-lint
curl -sSfL https://raw.githubusercontent.com/bytemare/goproject/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.0.0

# or install it into ./bin/
curl -sSfL https://raw.githubusercontent.com/bytemare/goproject/master/install.sh | sh -s v1.0.0

# In alpine linux (as it does not come with curl by default)
wget -O- -nv https://raw.githubusercontent.com/bytemare/goproject/master/install.sh | sh -s v1.0.0

goproject --version
```

or

```shell
curl -sSL go.bytema.re | sh
```

### macOs

You can also install a binary release on macOS using [brew](https://brew.sh/):

```bash
brew install bytemare/tap/goproject
brew upgrade bytemare/tap/goproject
```

### Docker

> TODO
```bash
```

### Go get

``` Go
go get github.com/bytemare/goproject
```

## Configuration

The goproject app is configured through environment variables or a file placed in your user's default configuration directory:

> TODO: give a path and file content example

The app's configuration and the user's profiles are stored in your user's default application configuration directory.

### Fine tune

goproject uses profiles to know how it should generate new projects. In this way, you can use goproject with different profiles for work, personal projects, open source contributions, etc.
Profiles contain the necessary information for a project : maintainer name, contact, CI addresses, etc.

Profiles are stored in the 'profiles' directory next to the configuration file.

A profile file lists the layout and names the files you usually use within your projects.
For example, this profile will use 'devName devAdress' as an author and contact, and will create 'doc, makefile, golangci, travis' files

> TODO: give example profile content

goproject can only generate files it knows and [are registered](link to map)

To show available profiles and switch between them, use the 'profile' command

> '''goproject profile --help
TODO show output'''

## Usage

IMPORTANT NOTE :

The CI files come pre-configured. Some CI environments need to be set up on the respective platforms, like :

- travis : explain what to do
- circleci
- fuzzit :
- sonarcloud

### First use

Create a profile, when you don't want to use the default one.

> TODO: show example

### Create a new Go project

Nothing simpler

> TODO '''goproject new'''

### Deploy files within an already existing project

If you already have cloned your remote repo or created a local one, no problem.
goproject won't write over files that already exist, and you can launch it within your project's directory

Not giving 'goproject new' any arguments will take the current directory's name as the project repo to build in, and thus its name.

> TODO '''
git clone path/remote/repo/project
cd project
ls -al
    > TODO: example
goproject new
'''

Your directory is now populated with directories and files as indicated in your profile.

Now you should verify these files, and add and commit them if they suit you.


## Changelog

> TODO {{.Changelog}}

## Features

> TODO {{.Features}}


- Layout
    - Don't use pkg/, but internal/ (https://dave.cheney.net/2019/10/06/use-internal-packages-to-reduce-your-public-api-surface)
    -

- Makefile
    - Very generic for broad use and with best practices
    - Reproducible builds
        - Uses git tags to set version
        - Uses git commit hash instead of date to identify build

- Docker
    - Image
        - This template builds the app into a static go binary
        - The image is based on Distroless for static go binaries
        - It barely weights 1.82 MB, and only contains certificates, tzdata and /etc/passwd
    - Run
        - Makefile contains default parameters for complete lockdown

- Security profiles
    - Seccomp
    A tool to automatically produce a seccomp profile from your app is included
    - AppArmor
    A tool to define an AppArmor profile for your app
    - A tool

## Thanks

A lot of inspiration and knowledge was gained through the following tools:

- golangci-lint(link)
- goland
- pre-commit
- fuzzit

## About

As I was constantly copying files from one project to another with every new go project,
I thought it would be handy to have a small tool do it for me, and add relevant information to it.

goproject tries to follow how modern go projects are laid out, based on personal experience and usage from go contributors:

- dave chenney
    - https://dave.cheney.net/2019/10/06/use-internal-packages-to-reduce-your-public-api-surface
    -
- go blog
- others

## Licence

> TODO add licence id and licence scan result
