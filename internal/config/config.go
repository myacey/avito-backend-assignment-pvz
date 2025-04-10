package config

import (
	"log"

	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web"
	"github.com/spf13/viper"
)

type AppConfig struct {
	PostgresHost     string `mapstructure:"POSTGRES_HOST"`
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresDB       string `mapstructure:"POSTGRES_DB"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`

	ServerCfg web.ServerConfig `mapstructure:"server"`
}

func LoadConfig() (config AppConfig, err error) {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	// viper.AutomaticEnv()

	viper.SetConfigFile("./configs/config.yaml")
	viper.MergeInConfig()

	viper.Unmarshal(&config)

	log.Println("config:", config)

	return
}
