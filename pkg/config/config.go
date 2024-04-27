package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database database `mapstructure:"database"`
	Redis    redis    `mapstructure:"redis"`
	Auth     auth     `mapstructure:"auth"`
}

type database struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DatabaseName string `mapstructure:"database_name"`
}

type redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type auth struct {
	Secret string `mapstructure:"secret"`
}

var config Config

func Init() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	return nil
}

func GetConfig() Config {
	return config
}
