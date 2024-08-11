package user

import (
	"context"
	modelUser "golang-rest-api/internal/model/user"
)

func (s UserService) UserProfile(ctx context.Context, userID string) (modelUser.UserProfileResp, error) {
	u, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return modelUser.UserProfileResp{}, err
	}

	return modelUser.UserProfileResp{
		ID:       u.ID,
		Name:     u.Name,
		Username: u.Username,
		Phone:    u.Phone,
	}, nil
}
