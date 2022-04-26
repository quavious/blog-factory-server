package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quavious/blog-factory-server/db"
)

func NewAdminMiddleware(repository *db.Repository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			username, isOK := c.Get("userID").(string)
			if !isOK {
				return c.JSON(http.StatusForbidden, &db.BadResponse{
					Status:  false,
					Message: "Invalid username.",
				})
			}
			var isAdmin bool
			row := repository.QueryRow("select isAdmin from users where username = ?", username)
			err := row.Scan(&isAdmin)
			if err != nil || !isAdmin {
				return c.JSON(http.StatusForbidden, &db.BadResponse{
					Status:  false,
					Message: "This user is unable to access.",
				})
			}
			return next(c)
		}
	}
}
