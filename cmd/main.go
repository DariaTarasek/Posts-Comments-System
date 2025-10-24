package main

import (
	"OzonTestTask/internal/config"
	"OzonTestTask/internal/storage/postgreSQL"
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Используются переменные окружения")
	}
	conf := config.GetConfig()
	fmt.Printf("Выбрано хранилище %s\n", conf.StorageType)
	if conf.StorageType == config.PostgresStorage {
		db, err := postgreSQL.NewDBConnection(conf.PostgresDSN)
		if err != nil {
			log.Fatalf("не удалось подключиться к БД: %v", err)
		}
		fmt.Println("Подключено хранилище postgres")
		defer db.Close()

	}
}
