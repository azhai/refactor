package cmd

import (
	"fmt"

	"gitea.com/azhai/refactor/config"
)

const (
	VERSION = "0.9.9"
)

var (
	settings *config.Settings
	verbose  bool // 详细输出
)

func Prepare(configFile string) *config.Settings {
	if configFile == "" {
		err := fmt.Errorf("need reverse file")
		panic(err)
	}
	var err error
	settings, err = config.ReadSettings(configFile)
	if err != nil {
		panic(err)
	}
	verbose = settings.Application.Debug
	return settings
}

func Verbose() bool {
	return verbose
}
