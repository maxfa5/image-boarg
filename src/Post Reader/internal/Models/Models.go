package Models

type Message struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	ChatID    int    `json:"chat_id"`
	CreatedAt string `json:"created_at"` // Change to String as Time is difficult
}
