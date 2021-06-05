package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type env struct {
	Port     string `mapstructure:"PORT"`
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBName   string `mapstructure:"DB_NAME"`
}

// ENV is a mapper to the environmets variables.
var ENV *env

// Load initializes the global ENV.
func Load() {
	ENV = &env{}
	env := strings.ToLower(os.Getenv("GO_ENVIRONMENT"))
	viper.SetConfigFile(fmt.Sprintf("%s.env", env))
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var vipErr viper.ConfigFileNotFoundError
		if ok := errors.As(err, &vipErr); ok {
			log.Fatalln(fmt.Errorf("config file not found. %w", err))
		} else {
			log.Fatalln(fmt.Errorf("unexpected error loading config file. %w", err))
		}
	}

	if err := viper.Unmarshal(ENV); err != nil {
		log.Fatalln(fmt.Errorf("failed to unmarshal config. %w", err))
	}
}
