package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/someengineering/fixctl/cmd"
)

var version = "dev"

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	cmd.SetVersion(version)
	if err := cmd.Execute(); err != nil {
		logrus.Errorln("Error:", err)
		os.Exit(1)
	}
}
