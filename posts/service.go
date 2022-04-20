package posts

import (
	"log"

	"github.com/quavious/blog-factory-server/config"
	"github.com/quavious/blog-factory-server/db"
)

type PostsService struct {
	config     *config.Config
	repository *db.Repository
}

func NewPostsService(config *config.Config, repository *db.Repository) *PostsService {
	return &PostsService{
		config:     config,
		repository: repository,
	}
}

func (service *PostsService) Posts(page int) []PostModel {
	res, err := service.repository.Query("select * from posts limit ?, ? order by created_at desc", (page)*10, (page-1)*10)
	if err != nil {
		return nil
	}
	posts := []PostModel{}
	for res.Next() {
		post := new(PostModel)
		res.Scan(post.ID, post.Title, post.Description, post.Content, post.CreatedAt, post.UpdatedAt, post.UserID)
		posts = append(posts, *post)
	}
	return posts
}

func (service *PostsService) Post(postID int) *PostModel {
	res := service.repository.QueryRow("select * from posts where id = ?", postID)
	post := new(PostModel)
	err := res.Scan(post.ID, post.Title, post.Description, post.Content, post.CreatedAt, post.UpdatedAt, post.UserID)
	if err != nil {
		return nil
	}
	return post
}

func (service *PostsService) Create(model *ModifyPostModel, userID string) bool {
	res, err := service.repository.Exec(`
	insert into posts (title, description, content, user_id) 
	values (?, ?, ?, ?)
	`, model.Title, model.Description, model.Content, userID)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	created, err := res.LastInsertId()
	if err != nil {
		log.Println("error: id is not integer.")
		return true
	}
	log.Println(created)
	return true
}

func (service *PostsService) Update(model *ModifyPostModel, postID int, userID string) bool {
	res, err := service.repository.Exec(`
	update posts set title = ?, description = ?, content = ? where id = ? and user_id = ?
	`, model.Title, model.Description, model.Content, postID, userID)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	updated, err := res.RowsAffected()
	if err != nil {
		log.Println("error: id is not integer.")
		return true
	}
	log.Println(updated)
	return true
}

func (service *PostsService) Delete(postID int, userID string) bool {
	res, err := service.repository.Exec(`
	delete from posts where id = ? and user_id = ?
	`, postID, userID)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	updated, err := res.RowsAffected()
	if err != nil {
		log.Println("error: id is not integer.")
		return true
	}
	log.Println(updated)
	return true
}
