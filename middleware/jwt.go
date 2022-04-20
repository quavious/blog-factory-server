package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/quavious/blog-factory-server/config"
	"github.com/quavious/blog-factory-server/db"
)

func NewJWTMiddleware(config *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		jwtAccessSecret, _ := config.GetJWTSecret()

		return func(c echo.Context) error {
			var err error
			header := c.Request().Header.Get("Authorization")
			if len(header) == 0 {
				err = errors.New("error: no authentication headers.")
				log.Println(err)
				return c.JSON(http.StatusForbidden, &db.BadResponse{
					Status:  false,
					Message: err.Error(),
				})
			}
			split := strings.Split(header, " ")
			if len(split) < 2 {
				err = errors.New("erorr: invalid token string")
				return c.JSON(http.StatusForbidden, &db.BadResponse{
					Status:  false,
					Message: err.Error(),
				})
			}
			accessToken, err := jwt.Parse(split[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("error: unexpected signing method")
				}
				return []byte(jwtAccessSecret), nil
			})
			if !accessToken.Valid || err != nil {
				err = errors.New("error: invalid tokens")
				return c.JSON(http.StatusForbidden, &db.BadResponse{
					Status:  false,
					Message: err.Error(),
				})
			}
			payload, ok := accessToken.Claims.(jwt.MapClaims)

			if !ok {
				err = errors.New("error: token parsing error")
				return c.JSON(http.StatusForbidden, &db.BadResponse{
					Status:  false,
					Message: err.Error(),
				})
			}
			c.Set("userID", payload["userId"])
			c.Set("email", payload["email"])
			return next(c)
		}
	}
}
