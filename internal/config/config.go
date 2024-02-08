package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	StorageDb  `yaml:"storage_db" env-required:"true"`
	HTTPServer `yaml:"http_server"`
	Auth       Auth
}

type StorageDb struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"8080"`
	DbName   string `yaml:"db_name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
type HTTPServer struct {
	Address       string        `yaml:"address" env-default:"localhost:8080"`
	Timeout       time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout   time.Duration `yaml:"idle_timeout" env-default:"60s"`
	AdminLogin    string        `yaml:"admin-login" env-default:"admin"`
	AdminPassword string        `yaml:"admin-password" env-required:"true" env:"ADMIN_PASSWORD"`
}
type Auth struct {
	SigningKey string `yaml:"signingKey" env-required:"true" env:"JWT_SIGNING_KEY"`
	Salt       string `yaml:"salt" env-required:"true" env:"SALT_FOR_HASH"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s ", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
