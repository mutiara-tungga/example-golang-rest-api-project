package postgres

import (
	"context"
	"errors"
	"fmt"

	"golang-rest-api/internal/entities"
	service "golang-rest-api/internal/service/user"
	"golang-rest-api/pkg/log"

	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) Create(ctx context.Context, u entities.User, actor string) error {
	const query = `INSERT INTO users (id, name, username, phone, password, created_by) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, u.ID, u.Name, u.Username, u.Phone, u.Password, actor)
	if err == nil {
		return nil
	}
	log.Error(ctx, "error create user", err)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "user_unique_username" {
		return service.ErrDuplicateUsername
	}
	return fmt.Errorf("postgres create user: %w", err)
}
