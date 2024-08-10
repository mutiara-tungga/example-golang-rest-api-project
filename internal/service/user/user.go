package user

import (
	"context"
	modelUser "golang-rest-api/internal/model/user"
	repoUser "golang-rest-api/internal/repository/user"
	"golang-rest-api/pkg/crypter"
	"golang-rest-api/pkg/database"

	"github.com/google/uuid"
)

type IUserService interface {
	CreateUser(ctx context.Context, req modelUser.CreateUserReq) (modelUser.CreateUserResp, error)
}

type UserServiceOption func(*UserService)

func WithUserRepo(userRepo repoUser.IUserRepo) UserServiceOption {
	return func(us *UserService) {
		us.userRepo = userRepo
	}
}

func WithTxHandler(db database.IPostgres) UserServiceOption {
	return func(us *UserService) {
		us.txHandler = database.NewTxHandler(db)
	}
}

type UserService struct {
	userRepo      repoUser.IUserRepo
	txHandler     database.TxHandler
	uuidGenerator func() string
	crypter       crypter.Crypter
}

func NewUserService(options ...UserServiceOption) UserService {
	res := &UserService{
		uuidGenerator: uuid.NewString,
		crypter:       crypter.New(),
	}

	for _, apply := range options {
		apply(res)
	}

	return *res
}
