package service

import (
	"context"
	"time"

	"golang-rest-api/internal/entities"
)

// UserRepository is the outbound persistence port.
type UserRepository interface {
	Create(ctx context.Context, user entities.User, actor string) error
	GetByID(ctx context.Context, id string) (entities.User, error)
	GetByUsername(ctx context.Context, username string) (entities.User, error)
}

// PasswordHasher is the outbound password port.
type PasswordHasher interface {
	Hash(ctx context.Context, password string) (string, error)
	Matches(ctx context.Context, password, hash string) bool
}

// TokenIssuer is the outbound token port.
type TokenIssuer interface {
	Issue(ctx context.Context, userID, username string) (TokenPair, error)
}

type TokenPair struct {
	AccessToken           string
	ExpiresAt             time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}
