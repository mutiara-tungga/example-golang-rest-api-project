package user

import "time"

const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

type InsertUser struct {
	ID       string
	Name     string
	Username string
	Phone    string
	Password string
	Actor    string
}

type User struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	Username string `db:"username"`
	Phone    string `db:"phone"`
	Password string `db:"password"`
}

type CreateUserReq struct {
	ID       string `json:"-"`
	Actor    string `json:"-"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type CreateUserResp struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
}

type UserLoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserLoginResp struct {
	AccessToken           string
	ExpiresAt             time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

type UserProfileResp struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
}
