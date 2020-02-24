package main

import (
	"github.com/bytemare/goproject/cmd"
)

func main() {
	/*logrus.WithFields(logrus.Fields{
	  "version": version,
	  "commit": commit,
	}).Info("Starting")*/

	cmd.Execute()
}
