// handlers содержит функции для работы с эндпоинтами.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ilyas-zal/UserManagmentTask/internal/database"
	"github.com/ilyas-zal/UserManagmentTask/internal/middlewares"
	"github.com/ilyas-zal/UserManagmentTask/internal/models"
)

// Auth производит аутентификацию пользователя и возвращает токен авторизации.
func Auth(w http.ResponseWriter, r *http.Request) {
	var user struct {
		UserID uint   `json:"user_id"`
		Email  string `json:"user_email"`
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неправильный запрос", http.StatusBadRequest)
		return
	}

	var dbUser models.User
	db := database.GetDB()
	db.First(&dbUser, user.UserID)
	if dbUser.ID == 0 || dbUser.Email != user.Email {
		http.Error(w, "Пользователь не найден или неверная почта", http.StatusUnauthorized)
		return
	}

	token, err := middlewares.GenerateToken(dbUser.ID)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// GetUserStatus возвращает информацию о пользователе по его ID.
func GetUserStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	db := database.GetDB()
	var user models.User
	db.First(&user, id)
	if user.ID == 0 {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// GetLeaderboard возвращает список пользователей, отсортированных по балансу JSON-формате.
func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	db := database.GetDB()
	db.Order("balance DESC").Find(&users)
	if len(users) == 0 {
		http.Error(w, "Список пользователей пуст", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// CompleteTask производит проверку и выполнение задания пользователем.
func CompleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	db := database.GetDB()
	var user models.User
	db.First(&user, id)
	if user.ID == 0 {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	var taskID struct {
		TaskID uint `json:"task_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&taskID)
	if err != nil {
		http.Error(w, "Неправильный запрос", http.StatusBadRequest)
		return
	}
	var task models.Task
	db.First(&task, taskID.TaskID)
	if task.ID == 0 {
		http.Error(w, "Задание не найдено", http.StatusNotFound)
		return
	}
	var userTask models.UserTask
	db.First(&userTask, "user_id = ? AND task_id = ?", user.ID, taskID.TaskID)
	if userTask.ID != 0 && userTask.Completed {
		http.Error(w, "Задание уже выполнено", http.StatusBadRequest)
		return
	}
	if userTask.ID == 0 {
		userTask.UserID = user.ID
		userTask.TaskID = taskID.TaskID
		userTask.Completed = true
		db.Create(&userTask)
	} else {
		userTask.Completed = true
		db.Save(&userTask)
	}
	user.Balance += task.Reward
	db.Save(&user)
	json.NewEncoder(w).Encode(user)
}

// SetReferrer устанавливает и награждает пригласившего(реферера) для пользователя.
func SetReferrer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	db := database.GetDB()
	var user models.User
	db.First(&user, id)
	if user.ID == 0 {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	var referrerID struct {
		ReferrerID uint `json:"referrer_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&referrerID)
	if err != nil {
		http.Error(w, "Неправильный запрос", http.StatusBadRequest)
		return
	}
	var referrer models.User
	db.First(&referrer, referrerID.ReferrerID)
	if referrer.ID == 0 {
		http.Error(w, "Реферер не найден", http.StatusNotFound)
		return
	}
	user.ReferrerID = &referrerID.ReferrerID
	db.Save(&user)
	referrer.Balance += 300
	db.Save(&referrer)
	json.NewEncoder(w).Encode(user)
}

// CreateUser производит регистрацию нового пользователя.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Name  string `json:"user_name"`
		Email string `json:"user_email"`
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неправильный запрос", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	var existingUser models.User
	db.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.ID != 0 {
		http.Error(w, "Пользователь с такой почтой уже существует", http.StatusBadRequest)
		return
	}

	newUser := models.User{
		Name:  user.Name,
		Email: user.Email,
	}
	db.Create(&newUser)
	json.NewEncoder(w).Encode(newUser)
}

// CreateTask производит создание новой задачи.
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task struct {
		Description string `json:"task_description"`
		Reward      int    `json:"task_reward"`
	}
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Неправильный запрос", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	newTask := models.Task{
		Description: task.Description,
		Reward:      task.Reward,
	}
	db.Create(&newTask)
	json.NewEncoder(w).Encode(newTask)
}
