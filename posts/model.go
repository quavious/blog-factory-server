package posts

import (
	"database/sql"
	"time"
)

type tagArray []string

type PostModel struct {
	ID          int            `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Content     string         `json:"content"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   sql.NullTime   `json:"updatedAt"`
	Username    string         `json:"username"`
	Tags        tagArray       `json:"tags,omitempty"`
	Comments    []CommentModel `json:"comments,omitempty"`
}

type CreatePostModel struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	Tags        []string `json:"tags"`
}

type UpdatePostModel struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

type CommentModel struct {
	ID        int          `json:"id"`
	Content   string       `json:"content"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt sql.NullTime `json:"updatedAt"`
	Username  string       `json:"username"`
}
