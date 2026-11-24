package service

import "context"

type Profile struct{ ID, Name, Username, Phone string }

func (s *Service) Profile(ctx context.Context, id string) (Profile, error) {
	u, err := s.users.GetByID(ctx, id)
	if err != nil {
		return Profile{}, err
	}
	return Profile{ID: u.ID, Name: u.Name, Username: u.Username, Phone: u.Phone}, nil
}
