package comments

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/quavious/blog-factory-server/config"
	"github.com/quavious/blog-factory-server/db"
)

type CommentsController struct {
	*echo.Echo
	config        *config.Config
	repository    *db.Repository
	jwtMiddleware *echo.MiddlewareFunc
}

func NewCommentsController(
	echo *echo.Echo,
	config *config.Config,
	repository *db.Repository,
	jwtMiddleware *echo.MiddlewareFunc,
) *CommentsController {
	return &CommentsController{
		Echo:          echo,
		config:        config,
		repository:    repository,
		jwtMiddleware: jwtMiddleware,
	}
}

func (controller *CommentsController) UseRoute() {
	commentsService := NewCommentsService(controller.config, controller.repository)
	controller.POST("/comments", func(c echo.Context) error {
		model := new(CreateCommentModel)
		err := c.Bind(model)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Invalid data form.",
			})
		}
		userID := c.Get("userID").(string)
		comments := commentsService.Create(model, userID)
		if comments == nil {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Creating new comment is failed.",
			})
		}
		return c.JSON(http.StatusCreated, echo.Map{
			"status":   true,
			"comments": comments,
			"message":  "New comment is created.",
		})
	}, *controller.jwtMiddleware)

	controller.PUT("/comments/:id", func(c echo.Context) error {
		model := new(UpdateCommentModel)
		err := c.Bind(model)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Invalid data form.",
			})
		}
		param := c.Param("id")
		id, err := strconv.Atoi(param)
		if err != nil || id < 1 {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Invalid comment id.",
			})
		}
		userID := c.Get("userID").(string)
		comments := commentsService.Update(model, id, userID)
		if comments == nil {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Deleting the comment is failed.",
			})
		}
		return c.JSON(http.StatusCreated, echo.Map{
			"status":   true,
			"comments": comments,
			"message":  "New comment is updated.",
		})
	}, *controller.jwtMiddleware)

	controller.DELETE("/comments/:id", func(c echo.Context) error {
		param := c.Param("id")
		id, err := strconv.Atoi(param)
		if err != nil || id < 1 {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Invalid comment id.",
			})
		}
		model := new(DeleteCommentModel)
		err = c.Bind(model)
		if err != nil {
			log.Println(err.Error())
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Invalid post id.",
			})
		}
		userID := c.Get("userID").(string)
		comments := commentsService.Delete(model, id, userID)
		if comments == nil {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Deleting the comment is failed.",
			})
		}
		return c.JSON(http.StatusOK, echo.Map{
			"status":   true,
			"comments": comments,
			"message":  "The comment is deleted.",
		})
	}, *controller.jwtMiddleware)
}
