package domain

import (
	"time"
)

type TodoList struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Date        *time.Time `json:"date"`
	CreatedAt   time.Time  `json:"created_at"`
}

type TodoListExtended struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Date        *time.Time `json:"date"`
	CreatedAt   time.Time  `json:"created_at"`
	UserId      int        `json:"user_id"`
}

type UserList struct {
	Id     int `json:"id"`
	UserId int `json:"user_id"`
	ListId int `json:"list_id"`
}

type TodoItem struct {
	Id                   int        `json:"id"`
	Title                string     `json:"title"`
	Description          string     `json:"description"`
	Done                 *bool      `json:"done"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	Note                 *string    `json:"note"`
	NotificationTime     *time.Time `json:"notification_time"`
	PredictedTimeToSpend *time.Time `json:"predicted_time_to_spend"`
	Priority             *int       `json:"priority"`
	IsDeleted            bool       `json:"is_deleted"`
	CategoryId           *int       `json:"category_id"`
	ActionMade           *bool      `json:"action_made"`
	ParentID             *int       `json:"parent_id"`
	Children             []TodoItem `json:"children"`
}

type Category struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IconName string `json:"icon_name"`
}

type ListItem struct {
	Id     int `json:"id"`
	ListId int `json:"list_id"`
	ItemId int `json:"item_id"`
}

type TodoItemGet struct {
	Id                   int           `json:"id"`
	Title                string        `json:"title"`
	Description          string        `json:"description"`
	Done                 *bool         `json:"done"`
	CreatedAt            time.Time     `json:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at"`
	Note                 *string       `json:"note"`
	NotificationTime     *time.Time    `json:"notification_time"`
	PredictedTimeToSpend *time.Time    `json:"predicted_time_to_spend"`
	Priority             *int          `json:"priority"`
	IsDeleted            bool          `json:"is_deleted"`
	Category             Category      `json:"category"`
	ParentID             *int          `json:"parent_id"`
	Children             []TodoItemGet `json:"children"`
}
