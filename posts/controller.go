package posts

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/quavious/blog-factory-server/config"
	"github.com/quavious/blog-factory-server/db"
)

type PostsController struct {
	*echo.Echo
	repository    *db.Repository
	config        *config.Config
	jwtMiddleware *echo.MiddlewareFunc
}

func NewPostsController(echo *echo.Echo, repository *db.Repository, config *config.Config, jwtMiddleware *echo.MiddlewareFunc) *PostsController {
	return &PostsController{
		Echo:          echo,
		repository:    repository,
		config:        config,
		jwtMiddleware: jwtMiddleware,
	}
}

func (controller *PostsController) UseRoute() {
	postsService := NewPostsService(controller.config, controller.repository)
	controller.POST("/post", func(c echo.Context) error {
		model := new(ModifyPostModel)
		err := c.Bind(model)
		userID, ok := c.Get("userID").(string)
		if err != nil || !ok {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Invalid data form.",
			})
		}
		isCreated := postsService.Create(model, userID)
		if !isCreated {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Creating new post is failed.",
			})
		}
		return c.JSON(http.StatusCreated, echo.Map{
			"status":  true,
			"message": "New post is created.",
		})
	}, *controller.jwtMiddleware)

	controller.GET("/posts/:page", func(c echo.Context) error {
		param := c.Param("page")
		page, err := strconv.Atoi(param)
		if err != nil || page < 1 {
			return c.JSON(http.StatusNotFound, &db.BadResponse{
				Status:  false,
				Message: "Invalid page.",
			})
		}
		posts := postsService.Posts(page)
		return c.JSON(http.StatusOK, echo.Map{
			"status": true,
			"posts":  posts,
		})
	})

	controller.GET("/posts/id/:id", func(c echo.Context) error {
		param := c.Param("id")
		id, err := strconv.Atoi(param)
		if err != nil || id < 1 {
			return c.JSON(http.StatusNotFound, &db.BadResponse{
				Status:  false,
				Message: "Invalid post id.",
			})
		}
		post := postsService.Post(id)
		return c.JSON(http.StatusOK, echo.Map{
			"status": true,
			"post":   post,
		})
	})

	controller.PUT("/post/id/:id", func(c echo.Context) error {
		model := new(ModifyPostModel)
		err := c.Bind(model)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Invalid data form.",
			})
		}
		param := c.Param("id")
		postID, err := strconv.Atoi(param)
		userID, ok := c.Get("userID").(string)
		if err != nil || !ok {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Invalid post id.",
			})
		}
		isUpdated := postsService.Update(model, postID, userID)
		if !isUpdated {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Updating the post is failed.",
			})
		}
		return c.JSON(http.StatusCreated, echo.Map{
			"status":  true,
			"message": "The post is updated.",
		})
	}, *controller.jwtMiddleware)

	controller.DELETE("/post/id/:id", func(c echo.Context) error {
		param := c.Param("id")
		postID, err := strconv.Atoi(param)
		userID, ok := c.Get("userID").(string)
		if err != nil || !ok {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Invalid post id.",
			})
		}
		isDeleted := postsService.Delete(postID, userID)
		if !isDeleted {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Deleting the post is failed.",
			})
		}
		return c.JSON(http.StatusCreated, echo.Map{
			"status":  true,
			"message": "The post is deleted.",
		})
	}, *controller.jwtMiddleware)

	controller.GET("/posts/tag/:tag/:page", func(c echo.Context) error {
		tag := c.Param("tag")
		param := c.Param("page")
		page, err := strconv.Atoi(param)
		if err != nil {
			return c.JSON(http.StatusNotFound, &db.BadResponse{
				Status:  false,
				Message: "Invalid page.",
			})
		}
		posts := postsService.PostsByTag(tag, page)
		return c.JSON(http.StatusOK, echo.Map{
			"status": true,
			"posts":  posts,
		})
	})

	controller.GET("/posts/search/:keyword/:page", func(c echo.Context) error {
		keyword := c.Param("keyword")
		param := c.Param("page")
		page, err := strconv.Atoi(param)
		if err != nil {
			return c.JSON(http.StatusNotFound, &db.BadResponse{
				Status:  false,
				Message: "Invalid page.",
			})
		}
		posts := postsService.PostsByKeyword(keyword, page)
		return c.JSON(http.StatusOK, echo.Map{
			"status": true,
			"posts":  posts,
		})
	})
}
