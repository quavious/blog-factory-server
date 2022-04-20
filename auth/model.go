package auth

import "github.com/golang-jwt/jwt"

type SignUpModel struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

type SignInModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JWTToken struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type JWTBody struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

type SendEmailModel struct {
	Email string `json:"email"`
}

type VerifyEmailModel struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type ModifyPasswordModel struct {
	CurrentPassword    string `json:"currentPassword"`
	NewPassword        string `json:"newPassword"`
	NewPasswordConfirm string `json:"newPasswordConfirm"`
}
