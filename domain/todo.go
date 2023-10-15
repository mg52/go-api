package domain

type Todo struct {
	UserID  int    `json:"user_id"`
	Content string `json:"content"`
}

type TodoRequest struct {
	Content string `json:"content"`
}
