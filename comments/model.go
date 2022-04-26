package comments

import "time"

type CreateCommentModel struct {
	Content string `json:"content"`
	PostID  int    `json:"postId"`
}

type UpdateCommentModel struct {
	Content string `json:"content"`
	PostID  int    `json:"postId"`
}

type DeleteCommentModel struct {
	PostID int `json:"postId"`
}

type CommentArray []CommentModel

type CommentModel struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Username  string    `json:"username"`
}
