package users

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Username     string `json:"username"`
	IsAdmin      bool   `json:"isAdmin"`
	RefreshToken string `json:"refreshToken"`
}

type UserAccount struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
}
