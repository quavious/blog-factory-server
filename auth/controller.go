package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quavious/blog-factory-server/config"
	"github.com/quavious/blog-factory-server/db"
	"github.com/quavious/blog-factory-server/mail"
)

type AuthController struct {
	*echo.Echo
	repository    *db.Repository
	config        *config.Config
	jwtMiddleware *echo.MiddlewareFunc
	mailClient    *mail.MailClient
}

func NewAuthController(
	echo *echo.Echo,
	config *config.Config,
	repository *db.Repository,
	jwtMiddleware *echo.MiddlewareFunc,
	mailClient *mail.MailClient,
) *AuthController {
	return &AuthController{
		repository:    repository,
		config:        config,
		Echo:          echo,
		jwtMiddleware: jwtMiddleware,
		mailClient:    mailClient,
	}
}

func (controller *AuthController) UseRoute() {
	authService := NewAuthService(controller.repository, controller.mailClient, controller.config)
	controller.POST("/auth/sign-up", func(c echo.Context) error {
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

	controller.POST("/auth/sign-in", func(c echo.Context) error {
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
			"status":  true,
			"message": "The confirmation email is sent.",
		})
	})

	controller.POST("/auth/email-token", func(c echo.Context) error {
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

	controller.POST("/auth/email-verification", func(c echo.Context) error {
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

	controller.GET("/auth/token-verification", func(c echo.Context) error {
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

	controller.POST("/auth/password-modification", func(c echo.Context) error {
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

	controller.POST("/auth/password-restoration", func(c echo.Context) error {
		model := new(RestorePasswordModel)
		err := c.Bind(model)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Invalid data form.",
			})
		}
		isOK := authService.RestorePassword(model)
		if !isOK {
			return c.JSON(http.StatusBadRequest, &db.BadResponse{
				Status:  false,
				Message: "Restoring password is failed.",
			})
		}
		return c.JSON(http.StatusCreated, echo.Map{
			"status":  true,
			"message": "Password is recovered.",
		})
	})

	// Not fully implemented
	// controller.POST("/auth/sign-in-confirmation", func(c echo.Context) error {
	// 	model := new(VerifyEmailModel)
	// 	err := c.Bind(model)
	// 	if err != nil {
	// 		log.Println(err.Error())
	// 		return c.JSON(http.StatusBadRequest, &db.BadResponse{
	// 			Status:  false,
	// 			Message: "Signing user is failed.",
	// 		})
	// 	}
	// 	isOK := authService.ConfirmSignIn(model)
	// 	if isOK {
	// 		return c.JSON(http.StatusBadRequest, &db.BadResponse{
	// 			Status:  false,
	// 			Message: "Signing user is failed.",
	// 		})
	// 	}
	// 	return c.JSON(http.StatusBadRequest, echo.Map{
	// 		"status":  true,
	// 		"accessToken": tokens.AccessToken,
	// 		"message": "A user logged in.",
	// 	})
	// })
}
