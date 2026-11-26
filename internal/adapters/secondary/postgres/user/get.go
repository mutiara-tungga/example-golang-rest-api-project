package postgres

import (
	"context"
	"fmt"

	"golang-rest-api/internal/entities"
	service "golang-rest-api/internal/service/user"
	"golang-rest-api/pkg/database"
	"golang-rest-api/pkg/log"
)

func (r *Repository) get(ctx context.Context, query, arg string) (entities.User, error) {
	var row userRow
	if err := r.db.Get(ctx, &row, query, arg); err == nil {
		return entities.User{ID: row.ID, Name: row.Name, Username: row.Username, Phone: row.Phone, Password: row.Password}, nil
	} else if err == database.RecordNotFound {
		return entities.User{}, service.ErrNotFound
	} else {
		log.Error(ctx, "error get user", err)
		return entities.User{}, fmt.Errorf("postgres get user: %w", err)
	}
}
