package user

import (
	"context"
	"golang-rest-api/internal/model"
	"golang-rest-api/internal/model/user"
	"golang-rest-api/pkg/database"
	pkgErr "golang-rest-api/pkg/error"
	"golang-rest-api/pkg/log"

	"github.com/jackc/pgx/v5"
)

type IUserRepo interface {
	CreateUserTx(ctx context.Context, tx pgx.Tx, args user.InsertUser) error
	GetUserByID(ctx context.Context, ID string) (user.User, error)
}

type UserRepo struct {
	db database.IPostgres
}

func NewUserRepo(db database.IPostgres) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r UserRepo) CreateUserTx(ctx context.Context, tx pgx.Tx, args user.InsertUser) error {
	query := `INSERT INTO user (name, username, phone, password, created_by, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := tx.Exec(
		ctx,
		query,
		args.Name,
		args.Username,
		args.Phone,
		args.Password,
		args.Actor,
	)

	if err != nil {
		log.Error(ctx, "error create user", err)
		return pkgErr.NewCustomErrWithOriginalErr(model.ErrorExecQuery, err)
	}

	return nil
}

func (r UserRepo) GetUserByID(ctx context.Context, ID string) (user.User, error) {
	query := `SELECT id, name, username, phone, password, created_by, created_at 
		FROM user
		WHERE id = $1`

	res := user.User{}
	err := r.db.Get(
		ctx,
		&res,
		query,
		ID,
	)

	if err != nil {
		log.Error(ctx, "error get user by id", err)
		return res, pkgErr.NewCustomErrWithOriginalErr(model.ErrorExecQuery, err)
	}

	return res, nil
}
