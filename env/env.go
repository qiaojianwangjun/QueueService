package env

import (
	"os"
)

const (
	Debug   = "debug"
	Test    = "test"
	Audit   = "audit"
	Release = "release"
)

var env string

func init() {
	env = os.Getenv("env")
	if env == "" {
		env = Debug
	}
}

func GetEnv() string {
	return env
}

func IsRelease() bool {
	return env == Release
}
