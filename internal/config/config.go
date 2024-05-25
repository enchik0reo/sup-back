package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string         `yaml:"env" env-required:"true"`
	TgToken    string         `env:"TG_API_BOT_TOKEN" env-required:"true"`
	TgBotHost  string         `yaml:"tgBothost"`
	TgAdmins   []string       `env:"TG_API_BOT_ADMINS" env-required:"true"`
	Salt       string         `env:"TG_SALT" env-required:"true"`
	CtxTimeout time.Duration  `yaml:"ctx_timeout"`
	Storage    Postgres       `yaml:"psql_storage"`
	Server     ApiServer      `yaml:"api_server"`
	Frontend   FrontendServer `yaml:"frontend"`
}

type Postgres struct {
	Driver   string `yaml:"db_driver"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type ApiServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"90s"`
}

type FrontendServer struct {
	Domains []string `yaml:"domains" env-required:"true"`
}

func MustLoad() *Config {
	cfg := &Config{}

	path, err := fetchConfigPath()
	if err != nil {
		panic(err)
	}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	cfg.Storage.Password = os.Getenv("POSTGRES_PASSWORD")
	if cfg.Storage.Password == "" {
		panic("postgress password is not specified in environment variables")
	}

	cfg.TgToken = os.Getenv("TG_API_BOT_TOKEN")
	if cfg.TgToken == "" {
		panic("telegram bot token is not specified in environment variables")
	}

	adminStr := os.Getenv("TG_API_BOT_ADMINS")

	cfg.TgAdmins = strings.Split(adminStr, ",")
	if cfg.TgAdmins == nil {
		panic("telegram bot admins is not specified in environment variables")
	}

	cfg.Salt = os.Getenv("TG_SALT")
	if cfg.Salt == "" {
		panic("salt is not specified in environment variables")
	}

	return cfg
}

func fetchConfigPath() (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", fmt.Errorf("can't load config: %v", err)
	}

	path := os.Getenv("CONFIG_PATH")

	if path == "" {
		return "", fmt.Errorf("config path is empty")
	}

	return path, nil
}
