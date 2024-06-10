package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
	Minio       `yaml:"minio"`
}

type HTTPServer struct {
	Addr        string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
}

type Minio struct {
	MinioAccessKeyID     string `yaml:"minio_access_key_id" env-required:"true"`
	MinioSecretAccessKey string `yaml:"minio_secret_access_key" env-required:"true"`
	MinioBucketName      string `yaml:"minio_bucket_name" env-required:"true"`
	MinioEndpoint        string `yaml:"minio_endpoint" env-required:"true"`
}

func MustLoad() *Config {
	var cfg Config

	configPath := fetchConfig()
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
		if configPath == "" {
			panic("CONFIG_PATH environment variable not set")
		}
	}

	// Проверка на сущестовоние файла yaml
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("failed to stat config path: %v", err)
	}

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	return &cfg
}

func fetchConfig() string {
	var configPath string

	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	return configPath
}
