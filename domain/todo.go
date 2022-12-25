package domain

type Todo struct {
	ID      int    `json:"id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}
