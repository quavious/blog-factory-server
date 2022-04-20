package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/quavious/blog-factory-server/config"
)

type Repository struct {
	*sql.DB
}

func NewRepository(config *config.Config) *Repository {
	uri := config.GetDBAddress()
	db, err := sql.Open("mysql", uri)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &Repository{
		DB: db,
	}
}
