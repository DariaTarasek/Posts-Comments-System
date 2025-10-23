package main

import (
	"OzonTestTask/internal/config"
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf(".env файл не найден")
	}
	conf := config.GetConfig()
	fmt.Printf("Сервер запущен на порту %s. Используется хранилище %s\n", conf.Port, conf.StorageType)

}
