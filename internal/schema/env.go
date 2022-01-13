package schema

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

var GlobalEnv = &env{}

type env struct {
	AppPort     int    `envconfig:"APP_PORT"`
	AppName     string `envconfig:"APP_NAME"`
	LogPath     string `envconfig:"LOG_PATH"`
	ConfigsPath string `envconfig:"CONFIGS_PATH"`
}

func init() {
	if err := GlobalEnv.load(); err != nil {
		panic("Load global env failed")
	}
}

func (e *env) load() error {
	if err := envconfig.Process("", e); err != nil {
		fmt.Printf("Failed to process env var: %s", err)
		return err
	}

	if e.AppPort == 0 {
		e.AppPort = 8080
	}
	if e.AppName == "" {
		e.AppName = "is-today-holiday"
	}
	if e.LogPath == "" {
		e.LogPath = fmt.Sprintf("/data/log/app/%s", e.AppName)
	}
	if e.ConfigsPath == "" {
		e.ConfigsPath = fmt.Sprintf("/data/app/%s/configs", e.AppName)
	}
	return nil
}
