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
	res, err := service.repository.Query(`
	select p.id, p.title, p.description, p.content, p.created_at, p.updated_at, u.username 
	from posts as p 
	join users as u on p.user_id = u.id
	order by created_at desc 
	limit ? 
	offset ?`, (page)*10, (page-1)*10)
	if err != nil {
		return nil
	}
	posts := []PostModel{}
	for res.Next() {
		post := new(PostModel)
		err := res.Scan(&post.ID, &post.Title, &post.Description, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.Username)
		if err != nil {
			log.Println(err.Error())
		}
		posts = append(posts, *post)
	}
	return posts
}

func (service *PostsService) Post(postID int) *PostModel {
	res := service.repository.QueryRow(`
	select p.id, p.title, p.description, p.content, p.created_at, p.updated_at, u.username
	from posts as p
	join users as u on p.user_id = u.id
	where p.id = ?;`, postID)
	post := new(PostModel)
	err := res.Scan(&post.ID, &post.Title, &post.Description, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.Username)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return post
}

func (service *PostsService) Create(model *CreatePostModel, userID string) bool {
	tags := []int{}
	for _, tag := range model.Tags {
		res, err := service.repository.Exec(`insert ignore into tags (tag) values (?)`, tag)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		inserted, err := res.LastInsertId()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if inserted == 0 {
			var tagID int
			row := service.repository.QueryRow(`select id from tags where tag = ?`, tag)
			err := row.Scan(&tagID)
			if err != nil {
				log.Println(err.Error())
			} else {
				tags = append(tags, tagID)
			}
		} else {
			tags = append(tags, int(inserted))
		}
	}
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

	for _, tagID := range tags {
		_, err = service.repository.Exec(`insert into posts_and_tags (post_id, tag_id) values (?, ?)`, int(created), tagID)
		if err != nil {
			log.Println(err.Error())
		}
	}
	return true
}

func (service *PostsService) Update(model *UpdatePostModel, postID int, userID string) bool {
	res, err := service.repository.Exec(`
	update posts 
	set title = ?, description = ?, content = ? where id = ? and user_id = ?
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
	delete from posts 
	where id = ? and user_id = ?
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

func (service *PostsService) PostsByTag(tag string, page int) []PostModel {
	var tagID int
	row := service.repository.QueryRow(`select id from tags where tag = ?`, tag)
	err := row.Scan(&tagID)
	if err != nil {
		return nil
	}
	res, err := service.repository.Query(`
	select p.id, p.title, p.description, p.content, p.created_at, p.updated_at, u.username
	from posts as p
	join users as u on p.user_id = u.id
	join posts_and_tags as pt on p.id = pt.post_id
	join tags as t on pt.tag_id = t.id
	where t.tag = ?
	order by created_at desc 
	limit ?
	offset ?`, tag, (page)*10, (page-1)*10)
	if err != nil {
		return nil
	}
	posts := []PostModel{}
	for res.Next() {
		post := new(PostModel)
		res.Scan(&post.ID, &post.Title, &post.Description, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.Username)
		posts = append(posts, *post)
	}
	return posts
}

func (service *PostsService) PostsByKeyword(keyword string, page int) []PostModel {
	res, err := service.repository.Query(`
	select p.id, p.title, p.description, p.content, p.created_at, p.updated_at, u.username
	from posts as p
	join users as u on p.user_id = u.id
	where p.title like ? 
	order by created_at desc 
	limit ? 
	offset ?`, "%"+keyword+"%", (page)*10, (page-1)*10)
	if err != nil {
		return nil
	}
	posts := []PostModel{}
	for res.Next() {
		post := new(PostModel)
		res.Scan(&post.ID, &post.Title, &post.Description, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.Username)
		posts = append(posts, *post)
	}
	return posts
}
