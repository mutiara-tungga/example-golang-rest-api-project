package user

import (
	"context"

	modelUser "golang-rest-api/internal/model/user"

	"github.com/jackc/pgx/v5"
)

func (s UserService) CreateUser(ctx context.Context, req modelUser.CreateUserReq) (modelUser.CreateUserResp, error) {
	hashPwdBytes, err := s.crypter.GenerateHash(ctx, req.Password)
	if err != nil {
		return modelUser.CreateUserResp{}, err
	}

	insertUserArgs := modelUser.InsertUser{
		ID:       s.uuidGenerator(),
		Name:     req.Name,
		Username: req.Username,
		Phone:    req.Phone,
		Password: string(hashPwdBytes),
		Actor:    req.Actor,
	}

	err = s.txHandler.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
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
