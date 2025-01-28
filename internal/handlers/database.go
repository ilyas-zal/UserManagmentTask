package handlers

import (
	"log"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("База данных не инициализирована")
	}
	return db
}

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
