package service

import (
	"context"

	"golang-rest-api/internal/entities"
)

type CreateUserCommand struct {
	Name, Username, Phone, Password, Actor string
}

type CreateUserResult struct {
	ID, Name, Username, Phone string
}

func (s *Service) CreateUser(ctx context.Context, cmd CreateUserCommand) (CreateUserResult, error) {
	hash, err := s.passwords.Hash(ctx, cmd.Password)
	if err != nil {
		return CreateUserResult{}, err
	}

	u := entities.User{ID: s.newID(), Name: cmd.Name, Username: cmd.Username, Phone: cmd.Phone, Password: hash}
	if err := s.users.Create(ctx, u, cmd.Actor); err != nil {
		return CreateUserResult{}, err
	}
	return CreateUserResult{ID: u.ID, Name: u.Name, Username: u.Username, Phone: u.Phone}, nil
}
