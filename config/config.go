package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var Version string

type Server struct {
	Host string `default:"0.0.0.0"`
	Port string `default:"8080" envconfig:"PORT"`
}

type Serial struct {
	Address string
}

type Configuration struct {
	Server Server
	Serial Serial
}

func LoadConfig(prefix, filename string) (Configuration, error) {
	if filename != "" {
		if fileEnv, err := godotenv.Read(filename); err == nil {
			for key, val := range fileEnv {
				os.Setenv(prefix+"_"+key, val)
			}
		}
	}

	var c Configuration
	err := envconfig.Process(prefix, &c)
	return c, err
}
