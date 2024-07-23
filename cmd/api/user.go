package main

import "golang-rest-api/internal/repository/user"

type UsersRepositories struct {
	UserRepo        user.UserRepo
	UserAddressRepo user.UserAddressRepo
}
