package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"time"
)

type App struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
}

type Server struct {
	Port            string        `envconfig:"PORT" default:"8080"`
	AdminPort       string        `envconfig:"ADMIN_PORT" default:"8081"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"1m"`
	MetricsPath     string        `envconfig:"METRICS_PATH" default:"/metrics"`
}

var Env = EnvConfig{}

type EnvConfig struct {
	App
	Server
}

func init() {
	err := envconfig.Process("", &Env)
	if err != nil {
		log.Fatalf("Fail to load env config : %v", err)
	}
}
