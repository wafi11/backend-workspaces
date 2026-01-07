package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database  DatabaseConfig `mapstructure:"database"`
	SecretKey SecretKey      `mapstructure:"secret_key"`
	Duration  Duration       `mapstructure:"duration"`
	Server    ServerConfig   `mapstructure:"server"`
	Docker    DockerConfig   `mapstructure:"docker"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

type DockerConfig struct {
	MaxContainers int    `mapstructure:"max_containers"`
	DefaultImage  string `mapstructure:"default_image"`
}

type SecretKey struct {
	JwtSecretKey string `mapstructure:"jwt_secret_key"`
}

type Duration struct {
	AccessToken int `mapstructure:"access_token_duration"`
	RefresToken int `mapstructure:"refresh_token_duration"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("/")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./")
	viper.AddConfigPath("/app")

	viper.AutomaticEnv()
	viper.BindEnv("database.host", "DATABASE_HOST")
	viper.BindEnv("database.port", "DATABASE_PORT")
	viper.BindEnv("database.user", "DATABASE_USER")
	viper.BindEnv("database.password", "DATABASE_PASSWORD")
	viper.BindEnv("database.name", "DATABASE_NAME")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
