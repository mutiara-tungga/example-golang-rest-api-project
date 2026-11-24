package postgres

import (
	"context"

	"golang-rest-api/internal/entities"
)

func (r *Repository) GetByID(ctx context.Context, id string) (entities.User, error) {
	return r.get(ctx, `SELECT id, name, username, phone, password FROM users WHERE id = $1 AND deleted_at IS NULL`, id)
}
