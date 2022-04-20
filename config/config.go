package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	dbName     string
	dbUser     string
	dbHost     string
	dbPort     string
	dbPassword string

	mailAddress   string
	mailApiKey    string
	mailSecretKey string

	jwtAccessSecret  string
	jwtRefreshSecret string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("no config files")
		return nil
	}

	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbPassword := os.Getenv("DB_PASSWORD")

	mailAddress := os.Getenv("MAIL_ADDRESS")
	mailApiKey := os.Getenv("MAIL_API_KEY")
	mailSecretKey := os.Getenv("MAIL_SECRET_KEY")

	jwtAccessSecret := os.Getenv("JWT_ACCESS_SECRET")
	jwtRefreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	return &Config{
		dbName:           dbName,
		dbUser:           dbUser,
		dbHost:           dbHost,
		dbPort:           dbPort,
		dbPassword:       dbPassword,
		mailAddress:      mailAddress,
		mailApiKey:       mailApiKey,
		mailSecretKey:    mailSecretKey,
		jwtAccessSecret:  jwtAccessSecret,
		jwtRefreshSecret: jwtRefreshSecret,
	}
}

func (config *Config) GetDBAddress() string {
	// root:pwd@tcp(127.0.0.1:3306)/testdb
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.dbUser, config.dbPassword, config.dbHost, config.dbPort, config.dbName)
}

func (config *Config) GetJWTSecret() (string, string) {
	return config.jwtAccessSecret, config.jwtRefreshSecret
}

func (config *Config) GetEmailKey() (string, string) {
	return config.mailApiKey, config.mailSecretKey
}

func (config *Config) GetEmailAddress() string {
	return config.mailAddress
}
