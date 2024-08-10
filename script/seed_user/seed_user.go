package main

import (
	"context"
	modelUser "golang-rest-api/internal/model/user"
	repoUser "golang-rest-api/internal/repository/user"
	"golang-rest-api/pkg/database"
	"golang-rest-api/pkg/log"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	log.InitLogger(log.LoggerMetaData{
		LogLevel:   "INFO",
		Service:    "script_seed_user",
		AppVersion: "v0.0.0",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users := []modelUser.InsertUser{
		// TODO: read from csv
	}

	db := database.NewPostgres(
	// TODO: add ENV
	// database.WithPostgresDBHost(""),
	// database.WithPostgresDBPort(""),
	// database.WithPostgresDBUser(""),
	// database.WithPostgresDBPassword(""),
	// database.WithPostgresDBName(""),
	)

	txHandler := database.NewTxHandler(db)
	userRepo := repoUser.NewUserRepo(db)

	err := txHandler.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		for _, u := range users {
			err := userRepo.CreateUserTx(ctx, tx, u)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(ctx, "error seed user", err)
	}

	log.Info(ctx, "success seeding user")
}
