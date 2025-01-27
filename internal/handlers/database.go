package handlers

import (
	"log"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func InitDB(dsn string) error {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Ошибка инициализации базы данных:", err)
		return err
	}
	return nil
}
