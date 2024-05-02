package main

import (
	"fmt"
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
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
