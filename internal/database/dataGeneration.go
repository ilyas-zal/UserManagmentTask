// database содержит функции для работы с базой данных.
package database

import (
	"fmt"
	"log"

	"github.com/ilyas-zal/UserManagmentTask/internal/models"
)

// SeedDB заполняет базу данных начальными данными: создает пользователей и задачи.
// Функция создает 10 пользователей и 5 задач, а затем сохраняет их в базе данных.
func SeedDB() error {
	db := GetDB()

	users := make([]*models.User, 10)
	for i := 0; i < 10; i++ {
		user := &models.User{
			Name:  fmt.Sprintf("user%d", i+1),
			Email: fmt.Sprintf("user%d@example.com", i+1),
		}
		users[i] = user
	}

	if err := db.CreateInBatches(users, 10).Error; err != nil {
		log.Printf("Ошибка создания пользователей: %v", err)
		return err
	}

	tasks := make([]*models.Task, 5)
	for i := 0; i < 5; i++ {
		task := &models.Task{
			Description: fmt.Sprintf("Задача %d", i+1),
			Reward:      100,
		}
		tasks[i] = task
	}

	if err := db.CreateInBatches(tasks, 5).Error; err != nil {
		log.Printf("Ошибка создания задач: %v", err)
		return err
	}

	return nil
}
