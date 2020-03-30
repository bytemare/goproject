// Package main calls the command
package main

import (
	"github.com/bytemare/goproject/internal/commands"
	"github.com/bytemare/goproject/internal/version"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.WithFields(logrus.Fields{
		"version": version.GetVersion(),
		"commit":  version.GetCommit(),
	}).Info("Starting")
	commands.Execute()
}
