package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/ilyas-zal/UserManagmentTask/internal/handlers"

	//"github.com/ilyas-zal/UserManagmentTask/internal/middlewares"
	"github.com/ilyas-zal/UserManagmentTask/internal/models"
	_ "github.com/lib/pq"
)

func main() {
	// Подключаемся к базе данных

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("Переменная окружения DB_DSN не задана")
	}

	err := handlers.InitDB(dsn)
	if err != nil {
		log.Fatal("Ошибка инициализации базы данных:", err)
	}

	// Создаем таблицы в базе данных
	if err := handlers.GetDB().AutoMigrate(&models.User{}, &models.Task{}, &models.UserTask{}); err != nil {
		log.Printf("Ошибка миграции базы данных: %v", err)
	}
	log.Println("Миграция базы данных завершена")

	// Вызываем функцию для создания пользователей и задач
	if err := handlers.SeedDB(); err != nil {
		log.Printf("Ошибка создания пользователей и задач: %v", err)
	}
	log.Println("Сгенерированы пользователи и задачи в базу данных")

	router := mux.NewRouter()
	router.HandleFunc("/auth", handlers.Auth).Methods("POST")
	router.HandleFunc("/users/{id}/status", handlers.GetUserStatus).Methods("GET")
	router.HandleFunc("/users/leaderboard", handlers.GetLeaderboard).Methods("GET")
	router.HandleFunc("/users/{id}/task/complete", handlers.CompleteTask).Methods("POST")
	router.HandleFunc("/users/{id}/referrer", handlers.SetReferrer).Methods("POST")

	//router.Use(middlewares.AuthMiddleware)

	log.Println("Запуск сервера на :8093")
	if err := http.ListenAndServe(":8093", router); err != nil {
		log.Printf("Ошибка при запуске сервера: %v", err)
	}
	log.Println("Сервер успешно запущен. Готовы к тестам")
}
