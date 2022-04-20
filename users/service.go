package users

import (
	"log"

	"github.com/quavious/blog-factory-server/config"
	"github.com/quavious/blog-factory-server/db"
)

type UsersService struct {
	config     *config.Config
	repository *db.Repository
}

func NewUsersService(config *config.Config, repository *db.Repository) *UsersService {
	return &UsersService{
		config:     config,
		repository: repository,
	}
}

func (service *UsersService) GetAccount(userID string) *UserAccount {
	model := new(UserAccount)
	row := service.repository.QueryRow("select id, email, username from users where id = ?", userID)
	err := row.Scan(&model.ID, &model.Email, &model.Username)
	if err != nil {
		log.Println("error: no users with this id.")
		return nil
	}
	return model
}
