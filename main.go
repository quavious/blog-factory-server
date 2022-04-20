package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quavious/blog-factory-server/auth"
	"github.com/quavious/blog-factory-server/config"
	"github.com/quavious/blog-factory-server/db"
	"github.com/quavious/blog-factory-server/mail"
	"github.com/quavious/blog-factory-server/middleware"
	"github.com/quavious/blog-factory-server/users"
)

func main() {
	config := config.NewConfig()
	if config == nil {
		return
	}
	repository := db.NewRepository(config)
	if repository == nil {
		return
	}
	mailClient := mail.NewMailClient(config)
	if mailClient == nil {
		return
	}
	jwtMiddleware := middleware.NewJWTMiddleware(config)
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Hello Go!",
		})
	})
	authController := auth.NewAuthController(e, config, repository, &jwtMiddleware, mailClient)
	usersController := users.NewUsersController(e, config, repository, &jwtMiddleware)

	authController.UseRoute()
	usersController.UseRoute()
	e.Logger.Fatal(e.Start("localhost:5000"))
}
