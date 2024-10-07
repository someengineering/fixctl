package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/someengineering/fixctl/cmd"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	if err := cmd.Execute(); err != nil {
		logrus.Errorln("Error:", err)
		os.Exit(1)
	}
}
