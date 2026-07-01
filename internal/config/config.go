package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB глобальная переменная для хранения соединения с базой данных.
var DB *gorm.DB

// ConnectDatabase устанавливает подключение к PostgreSQL и сохраняет его в глобальную переменную DB.
func ConnectDatabase() {
	// Настройка строки подключения к базе данных.
	dsn := "host=localhost user=myuser password=mypassword dbname=myapp port=5433 sslmode=disable"

	// Открываем соединение с базой данных через GORM.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}

	// Сохраняем подключение в глобальной переменной для использования в приложенииa.
	DB = db
	log.Println("Подключение к PostgreSQL установлено!")
}
