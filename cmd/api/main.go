package main

import (
	"context"
	"fmt"
	"golang-rest-api/config"
	handlerUser "golang-rest-api/internal/handler/user"
	modelUser "golang-rest-api/internal/model/user"
	repoUser "golang-rest-api/internal/repository/user"
	serviceUser "golang-rest-api/internal/service/user"
	"golang-rest-api/pkg/database"
	httpmiddleware "golang-rest-api/pkg/http_middleware"
	httpserver "golang-rest-api/pkg/http_server"
	"golang-rest-api/pkg/jwt"
	"golang-rest-api/pkg/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	config.LoadEnvConfig()
	log.InitLogger(log.LoggerMetaData{
		LogLevel:   "",
		Service:    config.Get().AppName,
		AppVersion: "v0.0.0",
	})

	jwtGenerator := jwt.NewJWTGenerator(
		jwt.JWTGeneratorWithIssuer(config.Get().AppName),
		jwt.JWTGeneratorWithSigningMethod(jwt.JWTSigningMethodNameRS256, config.Get().JWTRSAPrivateKey),
	)

	jwtValidator := jwt.NewJWTValidator(
		jwt.JWTValidatorWithSigningMethod(jwt.JWTSigningMethodNameRS256, config.Get().JWTRSAPublicKey),
		jwt.JWTValidatorWithValidIssuer(config.Get().AppName),
	)

	r := chi.NewRouter()
	r.Use(httpmiddleware.PanicRecoverer)

	posgresDB := database.NewPostgres(
		database.WithPostgresDBHost(config.Get().DatabaseHost),
		database.WithPostgresDBPort(config.Get().DatabasePort),
		database.WithPostgresDBUser(config.Get().DatabaseUser),
		database.WithPostgresDBPassword(config.Get().DatabasePass),
		database.WithPostgresDBName(config.Get().DatabaseName),
	)

	// repository
	userRepo := repoUser.NewUserRepo(posgresDB)

	// service
	userService := serviceUser.NewUserService(
		serviceUser.WithTxHandler(posgresDB),
		serviceUser.WithUserRepo(userRepo),
		serviceUser.WithJWTGenerator(jwtGenerator),
	)

	// handler
	userHandler := handlerUser.NewUserHandler(userService)

	// router
	r.Method(http.MethodPost, "/api/v1/user/login", httpserver.HandlerWithError(userHandler.Login))
	r.Group(func(r chi.Router) {
		r.Use(httpmiddleware.JWTAuthUser(jwtValidator, modelUser.AccessTokenCookieName))
		r.Method(http.MethodPost, "/api/v1/user", httpserver.HandlerWithError(userHandler.CreateUser))
		r.Method(http.MethodGet, "/api/v1/user/profile", httpserver.HandlerWithError(userHandler.UserProfile))
	})

	httpServer := http.Server{
		Addr:              ":" + config.Get().AppPort,
		Handler:           r,
		ReadTimeout:       5 * time.Minute,
		ReadHeaderTimeout: 5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
	}

	// Start HTTP server
	go func() {
		log.Info(context.Background(), fmt.Sprintf("Start Http Server at port %s", config.Get().AppPort))

		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(context.Background(), "Error HTTP server ListenAndServe: ", err)
		}
	}()

	// Capture SIGINT and SIGTERM signals for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Wait until we receive a shutdown signal
	<-signals

	// Create a context with a 30-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info(ctx, "Closing http server")

	// Attempt a graceful shutdown
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error(ctx, "Error server forced to shutdown: ", err)
	}

	log.Info(context.Background(), "Http server exiting gracefully")
}
