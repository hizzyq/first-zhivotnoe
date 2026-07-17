package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env    string `yaml:"env" env-default:"local"`
	Limits Limits `yaml:"limits"`
	Server Server `yaml:"server"`
	S3     S3     `yaml:"s3"`
	Broker Broker `yaml:"broker"`
}

type Limits struct {
	MaxSizeMB int `yaml:"max_size_mb"`
}

type Server struct {
	Port        string        `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type S3 struct {
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string `yaml:"access_key"`
	SecretKey  string `yaml:"secret_key"`
	BucketName string `yaml:"bucket_name"`
	UseSSL     string `yaml:"use_ssl"`
}

type Broker struct {
	Address   string `yaml:"address"`
	QueueName string `yaml:"queue_name"`
}

// MustLoad загружает конфиг или паникует
func MustLoad() *Config {
	// Ищем путь к конфигу
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// Если переменная не задана, берем локальный файл по умолчанию
		configPath = "config/local.yaml"
	}

	// Проверяем, существует ли файл
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	// Читаем YAML файл и заполняем структуру cfg
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
