package postgres

import (
	"context"

	"golang-rest-api/internal/entities"
)

func (r *Repository) GetByUsername(ctx context.Context, username string) (entities.User, error) {
	return r.get(ctx, `SELECT id, name, username, phone, password FROM users WHERE username = $1 AND deleted_at IS NULL`, username)
}
