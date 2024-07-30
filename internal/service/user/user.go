package user

import (
	"context"
	modelUser "golang-rest-api/internal/model/user"
	repoUser "golang-rest-api/internal/repository/user"
	"golang-rest-api/pkg/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
}

func NewUserService(options ...UserServiceOption) UserService {
	res := &UserService{
		uuidGenerator: uuid.NewString,
	}

	for _, apply := range options {
		apply(res)
	}

	return *res
}

func (s UserService) CreateUser(ctx context.Context, req modelUser.CreateUserReq) (modelUser.CreateUserResp, error) {
	insertUserArgs := modelUser.InsertUser{
		ID:       s.uuidGenerator(),
		Name:     req.Name,
		Username: req.Username,
		Phone:    req.Phone,
		Password: req.Password,
		Actor:    req.Actor,
	}

	err := s.txHandler.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		err := s.userRepo.CreateUserTx(ctx, tx, insertUserArgs)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return modelUser.CreateUserResp{}, err
	}

	return modelUser.CreateUserResp{
		ID:       insertUserArgs.ID,
		Name:     insertUserArgs.Name,
		Username: insertUserArgs.Username,
		Phone:    insertUserArgs.Phone,
	}, nil
}
