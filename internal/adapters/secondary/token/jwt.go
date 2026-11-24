package token

import (
	"context"

	service "golang-rest-api/internal/service/user"
	"golang-rest-api/pkg/jwt"
)

type JWT struct{ generator jwt.JWTGenerator }

func NewJWT(generator jwt.JWTGenerator) JWT { return JWT{generator: generator} }

func (j JWT) Issue(ctx context.Context, userID, username string) (service.TokenPair, error) {
	token, err := j.generator.GenerateJWT(ctx, jwt.User{ID: userID, Username: username})
	if err != nil {
		return service.TokenPair{}, err
	}
	return service.TokenPair{AccessToken: token.AccessToken, ExpiresAt: token.ExpiresAt, RefreshToken: token.RefreshToken, RefreshTokenExpiresAt: token.RefreshTokenExpiresAt}, nil
}
