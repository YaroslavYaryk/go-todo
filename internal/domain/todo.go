package domain

type TodoList struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoListExtended struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserId      int    `json:"user_id"`
}

type UserList struct {
	Id     int `json:"id"`
	UserId int `json:"user_id"`
	ListId int `json:"list_id"`
}

type TodoItem struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        *bool  `json:"done"`
}

type ListItem struct {
	Id     int `json:"id"`
	ListId int `json:"list_id"`
	ItemId int `json:"item_id"`
}
