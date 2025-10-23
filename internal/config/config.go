package config

import (
	"fmt"
	"log"
	"os"
)

type StorageType string

const (
	PostgresStorage StorageType = "postgres"
	InMemoryStorage StorageType = "memory"
)

type Config struct {
	Port        string
	StorageType StorageType
	PostgresDSN string
}

func NewConfig() *Config {
	storageType := getEnv("STORAGE_TYPE")
	conf := &Config{
		Port:        getEnv("PORT"),
		StorageType: StorageType(storageType),
	}

	if conf.StorageType == PostgresStorage {
		conf.PostgresDSN = getDSN()
	}

	return conf
}

func getEnv(key string) string {
	env := os.Getenv(key)
	if env == "" {
		log.Fatalf("не удалось получить переменную окружения %s", key)
	}
	return env
}

func getDSN() string {
	db := getEnv("POSTGRES_DB")
	user := getEnv("POSTGRES_USER")
	password := getEnv("POSTGRES_PASSWORD")
	host := getEnv("POSTGRES_HOST")
	port := getEnv("POSTGRES_PORT")
	ssl := os.Getenv("POSTGRES_SSLMODE")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?ssl=%s", user, password, host, port, db, ssl)
}

func GetConfig() *Config {
	conf := NewConfig()

	if conf.StorageType != PostgresStorage && conf.StorageType != InMemoryStorage {
		log.Fatalf("некорректный тип хранилища: %s", conf.StorageType)
	}

	return conf
}
