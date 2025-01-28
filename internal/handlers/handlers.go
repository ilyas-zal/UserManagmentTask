package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ilyas-zal/UserManagmentTask/internal/middlewares"
	"github.com/ilyas-zal/UserManagmentTask/internal/models"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неправильный запрос", http.StatusBadRequest)
		return
	}

	// Проверяем пользователя в базе данных
	var dbUser models.User
	db := GetDB()
	db.First(&dbUser, user.ID)
	if dbUser.ID == 0 {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	// Генерируем токен
	token, err := middlewares.GenerateToken(dbUser.ID)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func GetUserStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	db := GetDB()
	var user models.User
	db.First(&user, id)
	if user.ID == 0 {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	db := GetDB()
	db.Order("balance DESC").Find(&users)
	if len(users) == 0 {
		http.Error(w, "Список пользователей пуст", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func CompleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	db := GetDB()
	var user models.User
	db.First(&user, id)
	if user.ID == 0 {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	var task models.Task
	db.First(&task, 1) // Для примера берем первую задачу
	if task.ID == 0 {
		http.Error(w, "Задание не найдено", http.StatusNotFound)
		return
	}
	var userTask models.UserTask
	db.First(&userTask, "user_id = ? AND task_id = ?", user.ID, task.ID)
	if userTask.Completed {
		http.Error(w, "Задание уже выполнено", http.StatusBadRequest)
		return
	}
	userTask.Completed = true
	db.Save(&userTask)
	user.Balance += task.Reward
	db.Save(&user)
	json.NewEncoder(w).Encode(user)
}

func SetReferrer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	db := GetDB()
	var user models.User
	db.First(&user, id)
	if user.ID == 0 {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	var referrerID uint
	err := json.NewDecoder(r.Body).Decode(&referrerID)
	if err != nil {
		http.Error(w, "Неправильный запрос", http.StatusBadRequest)
		return
	}
	var referrer models.User
	db.First(&referrer, referrerID)
	if referrer.ID == 0 {
		http.Error(w, "Реферер не найден", http.StatusNotFound)
		return
	}
	user.ReferrerID = &referrerID
	db.Save(&user)
	json.NewEncoder(w).Encode(user)
}
