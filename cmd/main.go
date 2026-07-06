package main

import (
	"go-backend/internal/config"
	"go-backend/internal/handlers"
	"go-backend/internal/migrations"
	"go-backend/internal/models"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is not set in .env file")
	}

	config.ConnectDatabase()

	if err := migrations.RunMigrations(); err != nil {
		log.Fatal("Не удалось выполнить миграцию базы данных:", err)
	}

	if err := initData(); err != nil {
		log.Fatal("Не удалось инициализировать данные:", err)
	}

	r := gin.Default()

	// Передаём DB в роутер
	router, err := handlers.NewRouter(r, config.DB)
	if err != nil {
		log.Fatal("Не удалось создать роутер:", err)
	}
	router.RegisterRoutes()

	log.Println("🚀 Server starting on :8080")
	r.Run(":8080")
}

func initData() error {
	var count int64
	if err := config.DB.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		students := []models.User{
			{
				Fio:         "Иванов Иван",
				Group:       "AV",
				PhoneNumber: "+79001234567",
				Role:        "student",
				Password:    "",
			},
			{
				Fio:         "Петров Петр",
				Group:       "VM",
				PhoneNumber: "+79007654321",
				Role:        "student",
				Password:    "",
			},
			{
				Fio:         "Сидоров Сидор",
				Group:       "VM",
				PhoneNumber: "+79009999999",
				Role:        "student",
				Password:    "",
			},
			{
				Fio:         "Козлов Андрей",
				Group:       "AV",
				PhoneNumber: "+79008888888",
				Role:        "student",
				Password:    "",
			},
		}

		if err := config.DB.Create(&students).Error; err != nil {
			return err
		}
		log.Println("✅ Тестовые данные добавлены!")
	}
	return nil
}
