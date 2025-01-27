package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/ilyas-zal/UserManagmentTask/internal/handlers"
	"github.com/ilyas-zal/UserManagmentTask/internal/middlewares"
	"github.com/ilyas-zal/UserManagmentTask/internal/models"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

/* TODO:
инициализация функций
маршуты
запуск сервера
основная логика
*/

var db *gorm.DB

func main() {
	// Подключаемся к базе данных

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("Переменная окружения DB_DSN не задана")
	}

	// Ожидаем включения базы данных
	for {
		conn, err := net.Dial("tcp", "db:5432")
		if err != nil {
			log.Println("База данных не доступна, ожидаем...")
			time.Sleep(1 * time.Second)
			continue
		}
		conn.Close()
		break
	}

	err := handlers.InitDB(dsn)
	if err != nil {
		log.Fatal("Ошибка инициализации базы данных:", err)
	}

	// if db == nil {
	// 	log.Fatal("База данных не инициализирована после InitDB")
	// }

	// Создаем таблицы в базе данных
	if err := db.AutoMigrate(&models.User{}, &models.Task{}, &models.UserTask{}); err != nil {
		log.Fatalf("Ошибка миграции базы данных: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/auth", handlers.Auth).Methods("POST")
	router.HandleFunc("/users/{id}/status", handlers.GetUserStatus).Methods("GET")
	router.HandleFunc("/users/leaderboard", handlers.GetLeaderboard).Methods("GET")
	router.HandleFunc("/users/{id}/task/complete", handlers.CompleteTask).Methods("POST")
	router.HandleFunc("/users/{id}/referrer", handlers.SetReferrer).Methods("POST")

	router.Use(middlewares.AuthMiddleware)

	log.Println("Запуск сервера на :8093")
	if err := http.ListenAndServe(":8093", router); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	} else {
		log.Println("Сервер запущен")
	}
}
