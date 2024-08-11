package user

import (
	serviceUser "golang-rest-api/internal/service/user"
)

type UserHandler struct {
	userService serviceUser.IUserService
}

func NewUserHandler(
	userService serviceUser.IUserService,
) UserHandler {
	return UserHandler{
		userService: userService,
	}
}
