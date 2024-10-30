package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
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

func Init() *Environment {
	env := &Environment{}
	viper.SetConfigFile("config/local.env")
	viper.SetEnvPrefix("")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error when reading configuration file!\n", err)
		return env
	}

	// workaround for unmarshall not working with unbound env vars
	viper.SetDefault("SERVER_PORT", "5000")
	viper.SetDefault("ENV", "DEFAULT")
	viper.SetDefault("LOG_OUTPUT", "DEFAULT")
	viper.SetDefault("LOG_LEVEL", "DEFAULT")
	viper.SetDefault("DB_USER", "DEFAULT")
	viper.SetDefault("DB_PASS", "DEFAULT")
	viper.SetDefault("DB_HOST", "DEFAULT")
	viper.SetDefault("DB_PORT", "DEFAULT")
	viper.SetDefault("DB_NAME", "DEFAULT")
	viper.SetDefault("JWT_SECRET", "DEFAULT")
	err = viper.Unmarshal(env)
	if err != nil {
		log.Fatal("Error when parsing configuration", err)
	}
	fmt.Println(env)
	fmt.Println(env.Environment)
	return env
}
