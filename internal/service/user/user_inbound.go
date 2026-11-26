package service

import "context"

// UserUseCases is the inbound port used by primary adapters.
type UserUseCases interface {
	CreateUser(ctx context.Context, cmd CreateUserCommand) (CreateUserResult, error)
	Login(ctx context.Context, cmd LoginCommand) (TokenPair, error)
	Profile(ctx context.Context, id string) (Profile, error)
}
