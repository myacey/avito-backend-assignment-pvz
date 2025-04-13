package config

import (
	"log"

	"github.com/spf13/viper"

	pvzv1 "github.com/myacey/avito-backend-assignment-pvz/internal/grpc/pvz/v1"
	jwt_token "github.com/myacey/avito-backend-assignment-pvz/internal/pkg/jwt-token"
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

	TokenService jwt_token.TokenServiceConfig `mapstructure:"auth"`
}

func LoadConfig() (config AppConfig, err error) {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	viper.SetConfigFile("./configs/config.yaml")
	viper.MergeInConfig()

	viper.Unmarshal(&config)

	log.Println("config:", config)
	return
}
