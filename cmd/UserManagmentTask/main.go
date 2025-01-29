package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/ilyas-zal/UserManagmentTask/internal/database"
	"github.com/ilyas-zal/UserManagmentTask/internal/handlers"
	"github.com/ilyas-zal/UserManagmentTask/internal/middlewares"
	"github.com/ilyas-zal/UserManagmentTask/internal/models"
	_ "github.com/lib/pq"
)

func main() {

	// Берем данные для подключения к бд
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("Переменная окружения DB_DSN не задана")
	}
	// Инициализируем бд
	err := database.InitDB(dsn)
	if err != nil {
		log.Fatal("Ошибка инициализации базы данных:", err)
	}

	// Создаем таблицы в базе данных
	if err := database.GetDB().AutoMigrate(&models.User{}, &models.Task{}, &models.UserTask{}); err != nil {
		log.Printf("Ошибка миграции базы данных: %v", err)
	}
	log.Println("Миграция базы данных завершена")

	// Генерируем пользователей и задач, просто чтобы не пустые таблицы тестировать
	if err := database.SeedDB(); err != nil {
		log.Printf("Ошибка создания пользователей и задач: %v", err)
	}
	log.Println("Сгенерированы пользователи и задачи в базу данных")

	router := mux.NewRouter()
	router.HandleFunc("/auth", handlers.Auth).Methods("POST")
	router.HandleFunc("/users/{id}/status", middlewares.AuthMiddleware(handlers.GetUserStatus)).Methods("GET")
	router.HandleFunc("/users/leaderboard", middlewares.AuthMiddleware(handlers.GetLeaderboard)).Methods("GET")
	router.HandleFunc("/users/{id}/task/complete", middlewares.AuthMiddleware(handlers.CompleteTask)).Methods("POST")
	router.HandleFunc("/users/{id}/referrer", middlewares.AuthMiddleware(handlers.SetReferrer)).Methods("POST")
	// router.Use(middlewares.AuthMiddleware) как альтернатива, если эндпоинтов будет сильно больше, но нужно добавить исключения. Лучше так сделать
	//UPD: 2 новых эндпойнта для добавления нового пользователя и задач
	router.HandleFunc("/registration", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/newtask", handlers.CreateTask).Methods("POST")

	log.Println("Запуск сервера на :8093")
	if err := http.ListenAndServe(":8093", router); err != nil {
		log.Printf("Ошибка при запуске сервера: %v", err)
	}
	log.Println("Сервер успешно запущен. Готовы к тестам")
}
