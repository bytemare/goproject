# goproject
Set up a new go project with a default layout, with configuration files for your CI environment.

[![Build Status](https://travis-ci.com/bytemare/goproject.svg?branch=master)](https://travis-ci.com/bytemare/goproject)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=goproject&metric=alert_status)](https://sonarcloud.io/dashboard?id=goproject)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=goproject&metric=coverage)](https://sonarcloud.io/dashboard?id=goproject)
[![Go Report Card](https://goreportcard.com/badge/github.com/bytemare/goproject)](https://goreportcard.com/report/github.com/bytemare/goproject)
[![GolangCI](https://golangci.com/badges/github.com/bytemare/goproject.svg)](https://golangci.com/r/github.com/bytemare/goproject)
[![GoDoc](https://godoc.org/github.com/bytemare/goproject?status.svg)](https://godoc.org/github.com/bytemare/goproject)

Install the auto-updatable cli tool or simply use the template for a go project setup.

The templates contain basic CI configuration files, a Dockerfile, and some guidelines. Some example package functions and tests are also given.

## Installation or simple project structure

You can either install the tool to use it easily for all your next projects, with a simple command, or download the template.

## As an utility

To use it as a command line utility, that auto-updates, install it through go :

``` Go
go get github.com/bytemare/goproject
```

## Only the template

Through a shell download utility :

```shell
curl -sSL go.bytema.re | sh
```
