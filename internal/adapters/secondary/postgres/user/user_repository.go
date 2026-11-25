package postgres

import (
	service "golang-rest-api/internal/service/user"
	"golang-rest-api/pkg/database"
)

type Repository struct{ db database.IPostgres }

type userRow struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	Username string `db:"username"`
	Phone    string `db:"phone"`
	Password string `db:"password"`
}

func NewRepository(db database.IPostgres) *Repository { return &Repository{db: db} }

var _ service.UserRepository = (*Repository)(nil)
