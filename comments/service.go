package comments

import (
	"log"

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

func (service *CommentsService) Create(model *CreateCommentModel, userID string) bool {
	res, err := service.repository.Exec(`insert into comments (content, post_id, user_id) values (?, ?, ?)`, model.Content, model.PostID, userID)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	inserted, err := res.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return false
	}
	log.Println(inserted)
	return true
}

func (service *CommentsService) Update(model *UpdateCommentModel, id int, userID string) bool {
	res, err := service.repository.Exec(`update comments set content = ? where id = ? and user_id = ?`, model.Content, id, userID)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	updated, err := res.RowsAffected()
	if err != nil || updated == 0 {
		log.Println(err.Error())
		return false
	}
	return true
}

func (service *CommentsService) Delete(id int, userID string) bool {
	res, err := service.repository.Exec(`delete from comments where id = ? and user_id = ?`, id, userID)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	deleted, err := res.RowsAffected()
	if err != nil || deleted == 0 {
		log.Println(err.Error())
		return false
	}
	return true
}
