package main

import (
	"go-backend/internal/config"
	"go-backend/internal/handlers"
	"go-backend/internal/models"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()
	if err := config.DB.AutoMigrate(&models.Student{}); err != nil {
		log.Fatal("Не удалось выполнить миграцию базы данных:", err)
	}
	if err := initData(); err != nil {
		log.Fatal("Не удалось инициализировать данные:", err)
	}

	// Создаем новый экземпляр Gin с настройками по умолчанию.
	r := gin.Default()

	// Группируем маршруты под префиксом /api.
	api := r.Group("/api")
	{
		// Инициализируем обработчик студентов.
		studentHandler := handlers.NewStudentHandler()

		// Роуты для работы со студентами.
		api.GET("/students", studentHandler.GetAll)                                // Получить всех студентов.
		api.POST("/students", studentHandler.Create)                               // Создать нового студента.
		api.PUT("/students/:id", studentHandler.Update)                            // Обновить данные студента по ID.
		api.DELETE("/students/:id", studentHandler.Delete)                         // Удалить студента по ID.
		api.GET("/students/:id", studentHandler.GetByID)                           // Получить студента по ID.
		api.GET("/students/filter", studentHandler.FilterByGroup)                  // Фильтрация студентов по группе.
		api.GET("/students/filter-optional", studentHandler.FilterByGroupOptional) // Фильтрация с необязательным параметром.
		api.GET("/students/paginated", studentHandler.GetPaginated)                // Получить студентов постранично.
		api.GET("/students/search", studentHandler.Search)                         // Поиск студентов.
	}

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
