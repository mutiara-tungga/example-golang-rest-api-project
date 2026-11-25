package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang-rest-api/config"
	httpapi "golang-rest-api/internal/adapters/primary/httpapi/user"
	passwordadapter "golang-rest-api/internal/adapters/secondary/password"
	postgres "golang-rest-api/internal/adapters/secondary/postgres/user"
	tokenadapter "golang-rest-api/internal/adapters/secondary/token"
	service "golang-rest-api/internal/service/user"
	"golang-rest-api/pkg/database"
	"golang-rest-api/pkg/jwt"
	"golang-rest-api/pkg/log"

	_ "golang-rest-api/docs"

	"github.com/spf13/cobra"
)

func main() {
	// Capture SIGINT and SIGTERM signals for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	defer stop()

	if err := newRootCommand().ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:           "golang-project",
		Short:         "Golang Rest Project",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(&cobra.Command{
		Use:   "serve",
		Short: "Start the REST API server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runServer(cmd.Context())
		},
	})
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the application version",
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "golang-project v0.0.0")
		},
	})
	return root
}

func runServer(ctx context.Context) error {
	config.LoadEnvConfig()
	log.InitLogger(log.LoggerMetaData{LogLevel: "", Service: config.Get().AppName, AppVersion: "v0.0.0"})

	jwtGenerator := jwt.NewJWTGenerator(jwt.JWTGeneratorWithIssuer(config.Get().AppName), jwt.JWTGeneratorWithSigningMethod(jwt.JWTSigningMethodNameRS256, config.Get().JWTRSAPrivateKey))
	jwtValidator := jwt.NewJWTParser(jwt.JWTParserWithSigningMethod(jwt.JWTSigningMethodNameRS256, config.Get().JWTRSAPublicKey), jwt.JWTParserWithValidIssuer(config.Get().AppName))

	postgresDB := database.NewPostgres(
		database.WithPostgresDBHost(config.Get().DatabaseHost), database.WithPostgresDBPort(config.Get().DatabasePort),
		database.WithPostgresDBUser(config.Get().DatabaseUser), database.WithPostgresDBPassword(config.Get().DatabasePass),
		database.WithPostgresDBName(config.Get().DatabaseName),
	)
	defer postgresDB.Close()

	userRepo := postgres.NewRepository(postgresDB)
	userService := service.NewService(userRepo, passwordadapter.NewBcrypt(), tokenadapter.NewJWT(jwtGenerator))
	userHandler := httpapi.NewUserHandler(userService)
	router := httpapi.NewRouter(userHandler, jwtValidator)

	server := &http.Server{
		Addr:              ":" + config.Get().AppPort,
		Handler:           router,
		ReadTimeout:       5 * time.Minute,
		ReadHeaderTimeout: 5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
	}

	log.Info(ctx, fmt.Sprintf("Start Chi HTTP server at port %s", config.Get().AppPort))
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.ListenAndServe() }()

	select {
	case err := <-serverErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("HTTP server failed: %w", err)
	case <-ctx.Done():
		// Create a context with a 30-second timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		log.Info(shutdownCtx, "Closing HTTP server")
		// Attempt a graceful shutdown
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("HTTP server shutdown failed: %w", err)
		}
		if err := <-serverErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server failed while shutting down: %w", err)
		}
		log.Info(context.Background(), "HTTP server exiting gracefully")
		return nil
	}
}
