package service

import "github.com/google/uuid"

type Service struct {
	users     UserRepository
	passwords PasswordHasher
	tokens    TokenIssuer
	newID     func() string
}

func NewService(users UserRepository, passwords PasswordHasher, tokens TokenIssuer) *Service {
	return &Service{users: users, passwords: passwords, tokens: tokens, newID: uuid.NewString}
}
