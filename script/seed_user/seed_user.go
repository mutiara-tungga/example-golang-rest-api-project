package main

import (
	"context"
	"golang-rest-api/internal/adapters/secondary/password"
	postgres "golang-rest-api/internal/adapters/secondary/postgres/user"
	tokenadapter "golang-rest-api/internal/adapters/secondary/token"
	service "golang-rest-api/internal/service/user"
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

	users := []service.CreateUserCommand{
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
		database.WithPostgresDBHost(""),
		database.WithPostgresDBPort(""),
		database.WithPostgresDBUser(""),
		database.WithPostgresDBPassword(""),
		database.WithPostgresDBName(""),
	)

	userRepo := postgres.NewRepository(db)
	userService := service.NewService(userRepo, password.NewBcrypt(), tokenadapter.NewJWT(nil))

	for _, u := range users {
		if _, err := userService.CreateUser(ctx, u); err != nil {
			log.Fatal(ctx, "error seed user", err)
		}
	}

	log.Info(ctx, "success seeding user")
}
