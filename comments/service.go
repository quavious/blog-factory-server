package comments

import (
	"log"
	"time"

	"github.com/quavious/blog-factory-server/config"
	"github.com/quavious/blog-factory-server/db"
)

type CommentsService struct {
	config     *config.Config
	repository *db.Repository
}

func NewCommentsService(config *config.Config, repository *db.Repository) *CommentsService {
	return &CommentsService{
		config:     config,
		repository: repository,
	}
}

func (service *CommentsService) getComments(postID int) *CommentArray {
	rows, err := service.repository.Query(`
	select c.id, c.content, c.created_at, c.updated_at, u.username
	from comments as c
	join posts as p on c.post_id = p.id
	join users as u on c.user_id = u.id
	where c.post_id = ?
	order by c.created_at asc
	`, postID)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	comments := new(CommentArray)
	for rows.Next() {
		comment := new(CommentModel)
		err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.Username)
		if err != nil {
			log.Println(err.Error())
		} else {
			*comments = append(*comments, *comment)
		}
	}
	return comments
}

func (service *CommentsService) Create(model *CreateCommentModel, userID string) *CommentArray {
	createdAt := time.Now().UTC()
	res, err := service.repository.Exec(`insert into comments (content, post_id, created_at, updated_at, user_id) values (?, ?, ?, ?, ?)`, model.Content, model.PostID, createdAt, createdAt, userID)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	inserted, err := res.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	log.Println(inserted)
	comments := service.getComments(model.PostID)
	return comments
}

func (service *CommentsService) Update(model *UpdateCommentModel, id int, userID string) *CommentArray {
	updatedAt := time.Now().UTC()
	res, err := service.repository.Exec(`update comments set content = ?, updated_at = ? where id = ? and user_id = ?`, model.Content, updatedAt, id, userID)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	updated, err := res.RowsAffected()
	if err != nil || updated == 0 {
		log.Println(err.Error())
		return nil
	}
	comments := service.getComments(model.PostID)
	return comments
}

func (service *CommentsService) Delete(model *DeleteCommentModel, id int, userID string) *CommentArray {
	res, err := service.repository.Exec(`delete from comments where id = ? and user_id = ?`, id, userID)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	deleted, err := res.RowsAffected()
	if err != nil || deleted == 0 {
		log.Println(err.Error())
		return nil
	}
	comments := service.getComments(model.PostID)
	return comments
}
