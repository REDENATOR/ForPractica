package main

import (
	"go-backend/internal/config"
	"go-backend/internal/handlers"
	"go-backend/internal/migrations"
	"go-backend/internal/models"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()
	if err := migrations.RunMigrations(); err != nil {
		log.Fatal("Не удалось выполнить миграцию базы данных:", err)
	}
	if err := initData(); err != nil {
		log.Fatal("Не удалось инициализировать данные:", err)
	}

	// Создаем новый экземпляр Gin с настройками по умолчанию.
	r := gin.Default()

	router, err := handlers.NewRouter(r)
	if err != nil {
		log.Fatal("Не удалось создать роутер:", err)
	}
	router.RegisterRoutes()

	// Запускаем HTTP-сервер на порту 8080.
	r.Run(":8080")
}
func initData() error {
	var count int64
	if err := config.DB.Model(&models.Student{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		students := []models.Student{
			{Fio: "Иванов Иван", Group: "AV", PhoneNumber: "+79001234567"},
			{Fio: "Петров Петр", Group: "VM", PhoneNumber: "+79007654321"},
			{Fio: "Сидоров Сидор", Group: "VM", PhoneNumber: "+79009999999"},
			{Fio: "Козлов Андрей", Group: "AV", PhoneNumber: "+79008888888"},
		}
		if err := config.DB.Create(&students).Error; err != nil {
			return err
		}
		println("✅ Тестовые данные добавлены!")
	}
	return nil
}
