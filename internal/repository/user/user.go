package user

type IUserRepo interface{}

type UserRepo struct {
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}
