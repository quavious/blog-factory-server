package posts

import (
	"fmt"
	"log"
	"time"

	cm "github.com/quavious/blog-factory-server/comments"
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
	rows, err := service.repository.Query(`
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
	for rows.Next() {
		post := new(PostModel)
		err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.Username)
		if err != nil {
			log.Println(err.Error())
		}
		tags := service.Tags(post.ID)
		post.Tags = *tags
		posts = append(posts, *post)
	}
	return posts
}

func (service *PostsService) Post(postID int) (*PostModel, *cm.CommentArray) {
	res := service.repository.QueryRow(`
	select p.id, p.title, p.description, p.content, p.created_at, p.updated_at, u.username
	from posts as p
	join users as u on p.user_id = u.id
	where p.id = ?;`, postID)
	post := new(PostModel)
	err := res.Scan(&post.ID, &post.Title, &post.Description, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.Username)
	if err != nil {
		log.Println(err.Error())
		return nil, nil
	}
	fmt.Println(post.ID, post.CreatedAt)
	rows, err := service.repository.Query(`
	select c.id, c.content, c.created_at, c.updated_at, u.username
	from comments as c
	join posts as p on c.post_id = p.id
	join users as u on c.user_id = u.id
	where c.post_id = ?
	order by c.created_at asc
	`, post.ID)
	comments := new(cm.CommentArray)
	if err != nil {
		log.Println(err.Error())
	} else {
		for rows.Next() {
			comment := new(cm.CommentModel)
			err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.Username)
			if err != nil {
				log.Println(err.Error())
			} else {
				*comments = append(*comments, *comment)
			}
		}
	}
	post.Tags = *service.Tags(post.ID)
	return post, comments
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
	createdAt := time.Now().UTC()
	res, err := service.repository.Exec(`
	insert into posts (title, description, content, created_at, updated_at, user_id) 
	values (?, ?, ?, ?, ?, ?)
	`, model.Title, model.Description, model.Content, createdAt, createdAt, userID)
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
	updatedAt := time.Now().UTC()
	res, err := service.repository.Exec(`
	update posts 
	set title = ?, description = ?, content = ?, updated_at = ? where id = ? and user_id = ?
	`, model.Title, model.Description, model.Content, updatedAt, postID, userID)
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
		tags := service.Tags(post.ID)
		post.Tags = *tags
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
		tags := service.Tags(post.ID)
		post.Tags = *tags
		posts = append(posts, *post)
	}
	return posts
}

func (service *PostsService) Tags(postID int) *tagArray {
	rows, err := service.repository.Query(`
		select t.tag 
		from tags as t
		join posts_and_tags as pt on pt.tag_id = t.id
		join posts as p on p.id = pt.post_id
		where p.id = ?`, postID)
	if err != nil {
		return new(tagArray)
	} else {
		tags := new(tagArray)
		for rows.Next() {
			var tag string
			err := rows.Scan(&tag)
			if err != nil {
				log.Println(err.Error())
			} else {
				*tags = append(*tags, tag)
			}
		}
		return tags
	}
}
