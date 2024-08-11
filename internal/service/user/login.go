package user

import (
	"context"
	modelUser "golang-rest-api/internal/model/user"
	"golang-rest-api/pkg/jwt"
)

func (s UserService) UserLogin(ctx context.Context, req modelUser.UserLoginReq) (modelUser.UserLoginResp, error) {
	u, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return modelUser.UserLoginResp{}, err
	}

	passMatch := s.crypter.IsPWAndHashPWMatch(ctx, []byte(req.Password), []byte(u.Password))
	if !passMatch {
		return modelUser.UserLoginResp{}, modelUser.ErrorLoginErrorWrongPassword
	}

	jwtToken, err := s.jwtGenerator.GenerateJWT(ctx, jwt.User{
		ID:       u.ID,
		Username: u.Username,
	})
	if err != nil {
		return modelUser.UserLoginResp{}, err
	}

	return modelUser.UserLoginResp{
		AccessToken:           jwtToken.AccessToken,
		ExpiresAt:             jwtToken.ExpiresAt,
		RefreshToken:          jwtToken.RefreshToken,
		RefreshTokenExpiresAt: jwtToken.RefreshTokenExpiresAt,
	}, nil
}
