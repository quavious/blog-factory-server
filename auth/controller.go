package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quavious/blog-factory-server/config"
	"github.com/quavious/blog-factory-server/db"
	"github.com/quavious/blog-factory-server/mail"
)

type AuthController struct {
	repository    *db.Repository
	config        *config.Config
	echo          *echo.Echo
	jwtMiddleware *(func(echo.HandlerFunc) echo.HandlerFunc)
	mailClient    *mail.MailClient
}

func NewAuthController(
	echo *echo.Echo,
	config *config.Config,
	repository *db.Repository,
	jwtMiddleware *(func(echo.HandlerFunc) echo.HandlerFunc),
	mailClient *mail.MailClient,
) *AuthController {
	return &AuthController{
		repository:    repository,
		config:        config,
		echo:          echo,
		jwtMiddleware: jwtMiddleware,
		mailClient:    mailClient,
	}
}

func (controller *AuthController) Route() {
	authService := NewAuthService(controller.repository, controller.mailClient, controller.config)
	controller.echo.POST("/auth/sign-up", func(c echo.Context) error {
		model := new(SignUpModel)
		err := c.Bind(model)
		if err != nil || model.Password != model.PasswordConfirm {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Invalid data form.",
			})
		}
		isSigned := authService.SignUp(model)
		if !isSigned {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Bad request.",
			})
		}
		return c.JSON(http.StatusOK, echo.Map{
			"status":  true,
			"message": "New user signed up.",
			"email":   model.Email,
		})
	})

	controller.echo.POST("/auth/sign-in", func(c echo.Context) error {
		model := new(SignInModel)
		err := c.Bind(model)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Invalid data form.",
			})
		}
		tokens := authService.SignIn(model)
		if tokens == nil {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Signing user is failed.",
			})
		}
		c.SetCookie(&http.Cookie{
			Name:     "refreshToken",
			Value:    tokens.RefreshToken,
			HttpOnly: true,
			MaxAge:   7 * 24 * 60 * 60,
		})
		return c.JSON(http.StatusOK, echo.Map{
			"status":      true,
			"accessToken": tokens.AccessToken,
			"message":     "A user logged in.",
		})
	})

	controller.echo.POST("/auth/email-token", func(c echo.Context) error {
		model := new(SendEmailModel)
		err := c.Bind(model)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Invalid data form.",
			})
		}
		isOK := authService.SendEmail(model)
		if !isOK {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Sending new email is failed.",
			})
		}
		return c.JSON(http.StatusOK, echo.Map{
			"status":  true,
			"message": "The verification email is sent.",
		})
	})

	controller.echo.POST("/auth/email-verification", func(c echo.Context) error {
		model := new(VerifyEmailModel)
		err := c.Bind(model)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Invalid data form.",
			})
		}
		isOK := authService.VerifyEmail(model)
		if !isOK {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Email verification was failed.",
			})
		}
		return c.JSON(http.StatusOK, echo.Map{
			"status":  true,
			"message": "The email is now verified.",
		})
	})

	controller.echo.GET("/auth/token-verification", func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		if len(header) == 0 {
			return c.JSON(http.StatusForbidden, &db.BadResponse{
				Status:  false,
				Message: "no headers.",
			})
		}
		cookie, err := c.Cookie("refreshToken")
		if err != nil {
			return c.JSON(http.StatusForbidden, &db.BadResponse{
				Status:  false,
				Message: "no cookies.",
			})
		}
		accessToken := authService.VerifyJWTToken(header, cookie)
		if accessToken == nil {
			return c.JSON(http.StatusForbidden, &db.BadResponse{
				Status:  false,
				Message: "Token verification is failed.",
			})
		}
		return c.JSON(http.StatusOK, echo.Map{
			"status":      true,
			"accessToken": accessToken,
		})
	})

	controller.echo.POST("/auth/password-modification", func(c echo.Context) error {
		model := new(ModifyPasswordModel)
		err := c.Bind(model)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Invalid data form.",
			})
		}
		userID, ok := c.Get("userID").(string)
		if !ok {
			return c.JSON(http.StatusForbidden, &db.BadResponse{
				Status:  false,
				Message: "No user ids.",
			})
		}
		isOK := authService.ModifyPassword(userID, model)
		if !isOK {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Password mofidication is failed.",
			})
		}
		c.SetCookie(&http.Cookie{
			Name:     "refreshToken",
			Value:    "",
			HttpOnly: true,
			MaxAge:   -1,
		})
		return c.JSON(http.StatusOK, echo.Map{
			"status":  true,
			"message": "The password was changed.",
		})
	})
}
