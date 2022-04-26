package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/quavious/blog-factory-server/config"
	"github.com/quavious/blog-factory-server/db"
	"github.com/quavious/blog-factory-server/mail"
	"github.com/quavious/blog-factory-server/users"
	"github.com/quavious/blog-factory-server/utils"
)

type AuthService struct {
	repository *db.Repository
	mailClient *mail.MailClient
	config     *config.Config
}

func NewAuthService(repository *db.Repository, mailClient *mail.MailClient, config *config.Config) *AuthService {
	return &AuthService{
		repository: repository,
		mailClient: mailClient,
		config:     config,
	}
}

func (service *AuthService) SignUp(model *SignUpModel) bool {
	isVerified, err := service.repository.Query("select * from verified_emails where email = ?", model.Email)
	if err != nil {
		log.Println("error: this email is not verified")
		return false
	}
	index := 0
	for isVerified.Next() {
		index++
	}
	if index == 0 {
		log.Println("error: this email is not verified")
		return false
	}
	password, err := utils.Hash(model.Password)
	if err != nil {
		return false
	}
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Println("error: uuid is not created correctly.")
		return false
	}
	_, err = service.repository.Exec("insert into users (id, email, username, password) values (?, ?, ?, ?)", uuid.String(), model.Email, model.Username, password)
	if err != nil {
		log.Println("error: new user insertion is failed.")
		return false
	}
	return true
}

func (service *AuthService) SignIn(model *SignInModel) (*JWTToken, *users.User) {
	jwtAccessSecret, jwtRefreshSecret := service.config.GetJWTSecret()
	row := service.repository.QueryRow("select id, email, password, username, is_admin from users where email = ?", model.Email)
	user := new(users.User)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.IsAdmin)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	isOK, err := utils.Verify(model.Password, user.Password)
	if err != nil || !isOK {
		log.Println("error: password does not match.")
		return nil, nil
	}
	accessTokenclaims := JWTBody{
		UserID: user.ID,
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}
	refreshTokenclaims := JWTBody{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	}
	accessToken, err := (jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenclaims)).SignedString([]byte(jwtAccessSecret))
	if err != nil {
		return nil, nil
	}
	refreshToken, err := (jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenclaims)).SignedString([]byte(jwtRefreshSecret))
	if err != nil {
		return nil, nil
	}
	hashedRefreshToken, err := utils.Hash(refreshToken)
	if err != nil {
		return nil, nil
	}
	_, err = service.repository.Exec("update users set hashed_refresh_token = ? where id = ?", hashedRefreshToken, user.ID)
	if err != nil {
		return nil, nil
	}
	jwtToken := &JWTToken{AccessToken: accessToken, RefreshToken: refreshToken}
	return jwtToken, user
}

func (service *AuthService) SignOut() bool {
	return true
}

func (service *AuthService) SendEmail(model *SendEmailModel) bool {
	num := 0
	rows, err := service.repository.Query("select email from users where email = ?", model.Email)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	for rows.Next() {
		num++
	}
	if num > 0 {
		return false
	}
	emailToken := utils.NewEmailToken()
	expiredAt := time.Now().Add(time.Minute * 10)
	_, err = service.repository.Exec("insert into email_tokens (email, token, expired_at) values (?, ?, ?)", model.Email, emailToken, expiredAt)
	if err != nil {
		log.Println(err)
		return false
	}
	isOK := service.mailClient.SendToken(emailToken, model.Email)
	if !isOK {
		return false
	}
	return true
}

func (service *AuthService) VerifyEmail(model *VerifyEmailModel) bool {
	var expiredAt time.Time
	row := service.repository.QueryRow("select expired_at from email_tokens where email = ? and token = ? order by id desc", model.Email, model.Token)
	err := row.Scan(&expiredAt)
	if err != nil || !expiredAt.After(time.Now()) {
		return false
	}
	_, err = service.repository.Exec("insert ignore into verified_emails (email) values (?)", model.Email)
	if err != nil {
		return false
	}
	return true
}

func (service *AuthService) VerifyJWTToken(header string, cookie *http.Cookie) *string {
	split := strings.Split(header, " ")
	if len(split) < 2 {
		return nil
	}
	accessToken := split[1]
	refreshToken := cookie.Value

	tokens := &JWTToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	currentAccessToken := service.confirmJWTToken(tokens)
	return currentAccessToken
}

func (service *AuthService) ModifyPassword(userID string, model *ModifyPasswordModel) bool {
	if model.NewPassword != model.NewPasswordConfirm {
		return false
	}
	row := service.repository.QueryRow("select id, email, username, password from users where id = ?", userID)
	user := new(users.User)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		log.Println(err)
		return false
	}
	isOK, err := utils.Verify(model.CurrentPassword, user.Password)
	if err != nil || !isOK {
		return false
	}
	newHash, err := utils.Hash(model.NewPassword)
	if err != nil {
		log.Println(err)
		return false
	}
	_, err = service.repository.Exec("update users set password = ? where id = ?", newHash, user.ID)
	if err != nil {
		return false
	}
	return true
}

func (service *AuthService) RestorePassword(model *RestorePasswordModel) bool {
	if model.NewPassword != model.NewPasswordConfirm {
		return false
	}
	hash, err := utils.Hash(model.NewPassword)
	if err != nil {
		log.Println("error: hashing password is failed")
		return false
	}
	_, err = service.repository.Exec(`update users set password = ? where email = ?`, hash, model.Email)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func (service *AuthService) confirmJWTToken(tokens *JWTToken) *string {
	jwtAccessSecret, jwtRefreshSecret := service.config.GetJWTSecret()
	accessToken, err := jwt.Parse(tokens.AccessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error: unexpected signing method")
		}
		return []byte(jwtAccessSecret), nil
	})
	if accessToken.Valid {
		return &tokens.AccessToken
	}
	info, ok := err.(*jwt.ValidationError)
	if !ok {
		return nil
	}

	if info.Errors != jwt.ValidationErrorExpired {
		return nil
	}
	refreshToken, err := jwt.Parse(tokens.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error: unexpected signing method")
		}
		return []byte(jwtRefreshSecret), nil
	})
	if err != nil || !refreshToken.Valid {
		log.Println(err)
		return nil
	}
	claim, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("error: jwt map claims type casting")
		return nil
	}
	user := new(users.User)
	row := service.repository.QueryRow("select id, email, hashed_refresh_token from users where id = ?", claim["userId"])
	err = row.Scan(&user.ID, &user.Email, &user.RefreshToken)
	if err != nil {
		log.Println(err)
		return nil
	}
	isOK, err := utils.Verify(tokens.RefreshToken, user.RefreshToken)
	if !isOK || err != nil {
		log.Println("error: invalid refresh token")
		return nil
	}
	accessTokenClaims := JWTBody{
		UserID: user.ID,
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}
	newAccessToken, err := (jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)).SignedString([]byte(jwtAccessSecret))
	if err != nil {
		log.Println(err)
		return nil
	}
	return &newAccessToken
}

// Not fully implemented
// func (service *AuthService) ConfirmSignIn(model *VerifyEmailModel) bool {
// 	var expiredAt time.Time
// 	row := service.repository.QueryRow("select expired_at from email_tokens where email = ? and token = ? order by id desc", model.Email, model.Token)
// 	err := row.Scan(&expiredAt)
// 	if err != nil || !expiredAt.After(time.Now()) {
// 		return false
// 	}
// 	return true
// }
