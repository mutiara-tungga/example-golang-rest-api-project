package crypter

import (
	"context"
	pkgErr "golang-rest-api/pkg/error"
	"golang-rest-api/pkg/log"

	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -destination=mock/crypter.go -package=mock transport-service/pkg/crypter Crypter
type Crypter interface {
	GenerateHash(ctx context.Context, password string) ([]byte, error)
	IsPWAndHashPWMatch(ctx context.Context, password []byte, hashPass []byte) bool
}

func New() crypter {
	return crypter{}
}

type crypter struct{}

func (c crypter) GenerateHash(ctx context.Context, password string) ([]byte, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Error(ctx, "failed generate hash password", err)
		return nil, pkgErr.NewCustomErrWithOriginalErr(pkgErr.ErrFailedProcessPassword, err)
	}

	return passwordHash, nil
}

func (c crypter) IsPWAndHashPWMatch(ctx context.Context, password []byte, hashPass []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashPass, password)
	if err != nil {
		log.Error(ctx, "error compare password", err)
		return false
	}

	return true
}
