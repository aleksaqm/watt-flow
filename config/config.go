package config

import (
	"log"

	"github.com/spf13/viper"
)

type Environment struct {
	ServerPort  string `mapstructure:"SERVER_PORT"`
	Environment string `mapstructure:"ENV"`
	LogOutput   string `mapstructure:"LOG_OUTPUT"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
	DBUsername  string `mapstructure:"DB_USER"`
	DBPassword  string `mapstructure:"DB_PASS"`
	DBHost      string `mapstructure:"DB_HOST"`
	DBPort      string `mapstructure:"DB_PORT"`
	DBName      string `mapstructure:"DB_NAME"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
}

func Init() Environment {
	env := Environment{}
	viper.SetConfigFile("config/config.env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error when reading configuration file!")
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Error when parsing configuration file!")
	}
	return env
}
