package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/quavious/blog-factory-server/auth"
	"github.com/quavious/blog-factory-server/comments"
	"github.com/quavious/blog-factory-server/config"
	"github.com/quavious/blog-factory-server/db"
	"github.com/quavious/blog-factory-server/mail"
	md "github.com/quavious/blog-factory-server/middleware"
	"github.com/quavious/blog-factory-server/posts"
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
	defer repository.Close()
	mailClient := mail.NewMailClient(config)
	if mailClient == nil {
		return
	}
	jwtMiddleware := md.NewJWTMiddleware(config)
	corsMiddleware := md.NewCORSMiddleware()
	adminMiddleware := md.NewAdminMiddleware(repository)
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Hello Go!",
		})
	})
	e.Use(middleware.Recover())
	e.Use(*corsMiddleware)
	// e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	// 	return func(c echo.Context) error {
	// 		c.Response().Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	// 		return next(c)
	// 	}
	// })
	authController := auth.NewAuthController(e, config, repository, &jwtMiddleware, mailClient)
	usersController := users.NewUsersController(e, config, repository, &jwtMiddleware)
	postsController := posts.NewPostsController(e, config, repository, &jwtMiddleware, &adminMiddleware)
	commentsController := comments.NewCommentsController(e, config, repository, &jwtMiddleware)

	authController.UseRoute()
	usersController.UseRoute()
	postsController.UseRoute()
	commentsController.UseRoute()
	e.Logger.Fatal(e.Start("127.0.0.1:5000"))
}
