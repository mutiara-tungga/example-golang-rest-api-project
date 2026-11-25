package password

import (
	"context"

	service "golang-rest-api/internal/service/user"
	"golang-rest-api/pkg/crypter"
)

type Bcrypt struct{ crypter crypter.Crypter }

func NewBcrypt() *Bcrypt { return &Bcrypt{crypter: crypter.New()} }

func (b *Bcrypt) Hash(ctx context.Context, password string) (string, error) {
	hash, err := b.crypter.GenerateHash(ctx, password)
	return string(hash), err
}

func (b *Bcrypt) Matches(ctx context.Context, password, hash string) bool {
	return b.crypter.IsPWAndHashPWMatch(ctx, []byte(password), []byte(hash))
}

var _ service.PasswordHasher = (*Bcrypt)(nil)
