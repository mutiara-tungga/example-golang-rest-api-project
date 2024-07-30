package user

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
