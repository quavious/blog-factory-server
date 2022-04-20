package comments

type CreateCommentModel struct {
	Content string `json:"content"`
	PostID  int    `json:"postId"`
}

type UpdateCommentModel struct {
	Content string `json:"content"`
}
