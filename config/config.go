package config

var Version string = "dev"

func GetUserAgent() string {
	return "fixctl-" + Version
}
