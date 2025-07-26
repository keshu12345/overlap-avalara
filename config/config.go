package config

import (
	"os"
	"path"

	"github.com/keshu12345/overlap-avalara/toolkit"
	logger "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

const (
	serverYML = "server.yml"
)

// NewFxModule returns the fx.Option that builds the *Configuration struct
// that could be later used by other fx modules.
func NewFxModule(configDirPath string, overridePath string) fx.Option {
	return fx.Provide(
		func() (*Configuration, error) {
			var conf Configuration
			if len(configDirPath) == 0 {
				logger.Info("trying env config path ")
				configDirPath = os.Getenv("CONFIG_PATH")
			}
			logger.Info("Using config path ", path.Join(configDirPath, serverYML))
			err := toolkit.NewConfig(&conf, path.Join(configDirPath, serverYML), overridePath)
			return &conf, err
		},
	)
}

type Configuration struct {
	EnvironmentName string
	Server          Server       `mapstructure:"server"`
	Logger          LoggerConfig `mapstructure:"logger"`
}

type Server struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

type LoggerConfig struct {
	Base         string `yaml:"base"`         // e.g., "logrus"
	Level        string `yaml:"level"`        // e.g., "info", "debug"
	Format       string `yaml:"format"`       // e.g., "json", "text"
	ReportCaller bool   `yaml:"reportCaller"` // true to include caller info
	Enabled      bool   `yaml:"enabled"`      // enable logging output
	MaxSize      int    `yaml:"maxSize"`      // max size in MB before rotation
	MaxAge       int    `yaml:"maxAge"`       // days to retain old log files
	MaxBackups   int    `yaml:"maxBackups"`   // number of old files to retain
	LocalTime    bool   `yaml:"localTime"`    // use local time for timestamps
	Compress     bool   `yaml:"compress"`     // compress rotated logs
	LogDir       string `yaml:"logDir"`       // base directory for logs
}
