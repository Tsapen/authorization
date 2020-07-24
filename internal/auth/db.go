package auth

import "time"

// User contains registration info.
type User struct {
	Email       string `json:"email"`
	Login       string `json:"login"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

// Auth contains authentication info.
type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// TokenDetails contains token details.
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    time.Time
	RtExpires    time.Time
}

// DB is a database interface.
type DB interface {
	CreateUser(User) error
	CheckUser(Auth) error
	CreateToken(TokenDetails) error
}
