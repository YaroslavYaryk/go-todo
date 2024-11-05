package domain

type User struct {
	Id           int64   `json:"id" db:"id"`
	Name         string  `json:"name" binding:"required"`
	Username     string  `json:"username" binding:"required"`
	Password     string  `json:"password" binding:"required"`
	Theme        string  `json:"theme"`
	Timezone     string  `json:"timezone"`
	IsPaidMember bool    `json:"is_paid_member"`
	RateId       int     `json:"rate_id"`
	TaskColor    *string `json:"task_color"` // New field added

}

type UserGet struct {
	Id           int64   `json:"id" db:"id"`
	Name         string  `json:"name" binding:"required"`
	Username     string  `json:"username" binding:"required"`
	Theme        string  `json:"theme"`
	Timezone     *string `json:"timezone"`
	IsPaidMember bool    `json:"is_paid_member"`
	RateId       *int    `json:"rate_id"`
	TaskColor    *string `json:"task_color"` // New field added
}

type Rate struct {
	Id             int64  `json:"id" db:"id"`
	Name           string `json:"name"`
	TasksCompleted int    `json:"tasks_completed"`
}
