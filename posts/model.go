package posts

import "time"

type PostModel struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	UserID      string    `json:"userId"`
}

type ModifyPostModel struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}
