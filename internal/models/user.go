// models содержит структуру данных пользователей и задач.
package models

type User struct {
	ID         uint   `gorm:"primaryKey" json:"user_id"`
	Name       string `gorm:"unique;not null" json:"user_name"`
	Email      string `gorm:"unique;not null" json:"user_email"`
	Balance    int    `json:"user_balance"`
	ReferrerID *uint  `json:"referrer_id"`
	Referrer   *User  `json:"referrer" gorm:"foreignKey:ReferrerID;references:ID"`
}

type Task struct {
	ID          uint   `gorm:"primaryKey" json:"task_id"`
	Description string `json:"task_description"`
	Reward      int    `json:"task_reward"`
}

type UserTask struct {
	ID        uint  `gorm:"primaryKey" json:"user_task_id"`
	UserID    uint  `json:"user_id" gorm:"foreignKey:UserID;references:ID"`
	TaskID    uint  `json:"task_id" gorm:"foreignKey:TaskID;references:ID"`
	Completed bool  `json:"task_completed"`
	User      *User `json:"user"`
	Task      *Task `json:"task"`
}

/*
type User struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `gorm:"unique;not null" json:"name"`
	Email      string `gorm:"unique;not null" json:"email"`
	Balance    int    `json:"money"`
	ReferrerID *uint  `json:"referrer_id"`                                         // ID реферера, может быть nil
	Referrer   *User  `json:"referrer" gorm:"foreignKey:ReferrerID;references:ID"` // Связь с реферером
}

type Task struct {
	ID          uint   `gorm:"primaryKey" json:"taskid"`
	Description string `json:"description"`
	Reward      int    `json:"reward"`
}

type UserTask struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	UserID    uint `json:"user_id"`
	TaskID    uint `json:"task_id"`
	Completed bool `json:"completed"`
}

type Referral struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	UserID     uint `json:"user_id"`
	ReferrerID uint `json:"referrer_id"`
	Reward     int  `json:"reward"`
}
*/
