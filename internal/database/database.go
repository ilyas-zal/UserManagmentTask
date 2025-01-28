// database содержит функции для работы с базой данных.
package database

import (
	"log"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var db *gorm.DB

// GetDB возвращает экземпляр базы данных GORM.
// Если база данных не инициализирована, функция завершает выполнение с ошибкой.
func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("База данных не инициализирована")
	}
	return db
}

// InitDB инициализирует подключение к базе данных с использованием предоставленной строки подключения (DSN).
func InitDB(dsn string) error {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Ошибка инициализации базы данных: %v", err)
		return err
	}
	log.Println("База данных инициализирована успешно")
	return nil
}
