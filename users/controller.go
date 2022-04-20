package users

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quavious/blog-factory-server/config"
	"github.com/quavious/blog-factory-server/db"
)

type UsersController struct {
	*echo.Echo
	config        *config.Config
	repository    *db.Repository
	jwtMiddleware *(func(echo.HandlerFunc) echo.HandlerFunc)
}

func NewUsersController(echo *echo.Echo, config *config.Config, repository *db.Repository, jwtMiddleware *(func(echo.HandlerFunc) echo.HandlerFunc)) *UsersController {
	return &UsersController{Echo: echo, config: config, repository: repository, jwtMiddleware: jwtMiddleware}
}

func (controller *UsersController) Route() {
	userService := NewUsersService(controller.config, controller.repository)
	controller.GET("/users/account", func(c echo.Context) error {
		userID := c.Get("userID").(string)
		account := userService.GetAccount(userID)
		if account == nil {
			return c.JSON(http.StatusForbidden, &db.BadResponse{
				Status:  false,
				Message: "No users.",
			})
		}
		return c.JSON(http.StatusOK, echo.Map{
			"status":  true,
			"account": account,
		})
	}, *controller.jwtMiddleware)
}
