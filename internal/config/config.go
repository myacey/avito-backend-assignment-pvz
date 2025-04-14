package config

import (
	"log"

	"github.com/spf13/viper"

	pvzv1 "github.com/myacey/avito-backend-assignment-pvz/internal/grpc/pvz/v1"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/jwttoken"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web"
)

type AppConfig struct {
	PostgresHost     string `mapstructure:"POSTGRES_HOST"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT"`
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresDB       string `mapstructure:"POSTGRES_DB"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`

	HTTPServerCfg web.ServerConfig `mapstructure:"httpserver"`
	GRPCServerCfg pvzv1.Config     `mapstructure:"grpcserver"`

	TokenService jwttoken.TokenServiceConfig `mapstructure:"auth"`
}

func LoadConfig(cfgPath string) (config AppConfig, err error) {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()

	viper.SetConfigFile(cfgPath)
	viper.MergeInConfig()

	viper.Unmarshal(&config)

	log.Println("config:", config)

	return
}
