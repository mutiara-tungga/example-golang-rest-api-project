package service

import "context"

type LoginCommand struct{ Username, Password string }

func (s *Service) Login(ctx context.Context, cmd LoginCommand) (TokenPair, error) {
	u, err := s.users.GetByUsername(ctx, cmd.Username)
	if err != nil {
		return TokenPair{}, err
	}
	if !s.passwords.Matches(ctx, cmd.Password, u.Password) {
		return TokenPair{}, ErrWrongPassword
	}
	return s.tokens.Issue(ctx, u.ID, u.Username)
}
