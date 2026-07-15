package config

import (
	"fmt"
	"net/url"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Minio    MinioConfig
}

type AppConfig struct {
	Host string `env:"HOST" env-default:"0.0.0.0"`
	Port string `env:"PORT" env-default:"8080"`
}

type PostgresConfig struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
	SSLMode  string `env:"DB_SSL_MODE" env-default:"disable"`
}

type MinioConfig struct {
	Endpoint        string `env:"MINIO_ENDPOINT"`
	AccessKeyID     string `env:"ACCESS_KEY_ID"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY"`
	UseSSL          bool   `env:"USE_SSL" env-default:"false"`
	BucketName      string `env:"BUCKET_NAME"`
}

func New() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *PostgresConfig) DSN() string {
	return (&url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:   cfg.Name,
		RawQuery: url.Values{
			"sslmode": {cfg.SSLMode},
		}.Encode(),
	}).String()
}
