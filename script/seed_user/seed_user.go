package main

import (
	"context"
	postgres "golang-rest-api/internal/adapters/secondary/postgres/user"
	"golang-rest-api/internal/entities"
	"golang-rest-api/pkg/database"
	"golang-rest-api/pkg/log"
	"time"
)

func main() {
	log.InitLogger(log.LoggerMetaData{
		LogLevel:   "INFO",
		Service:    "script_seed_user",
		AppVersion: "v0.0.0",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users := []entities.User{
		{
			Name:     "Admin User",
			Username: "admin",
			Phone:    "",
			Password: "password",
		},
		{
			Name:     "Regular User",
			Username: "user",
			Phone:    "",
			Password: "password",
		},
	}

	db := database.NewPostgres(
	// TODO: add ENV
	// database.WithPostgresDBHost(""),
	// database.WithPostgresDBPort(""),
	// database.WithPostgresDBUser(""),
	// database.WithPostgresDBPassword(""),
	// database.WithPostgresDBName(""),
	)

	userRepo := postgres.NewRepository(db)

	for _, u := range users {
		if err := userRepo.Create(ctx, u, "seed"); err != nil {
			log.Fatal(ctx, "error seed user", err)
		}
	}

	log.Info(ctx, "success seeding user")
}
