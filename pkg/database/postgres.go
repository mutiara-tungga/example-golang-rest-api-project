package database

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ IPostgres = (*Postgres)(nil)

type PostgresConfigOption func(*PostgresConfig)

type PostgresConfig struct {
	User           string `validate:"required"`
	Password       string `validate:"required"`
	Host           string `validate:"required"`
	Port           string `validate:"required"`
	DatabaseName   string `validate:"required"`
	OptionalConfig map[string]string
}

func WithPostgresDBUser(user string) PostgresConfigOption {
	return func(pc *PostgresConfig) {
		pc.User = user
	}
}

func WithPostgresDBPassword(password string) PostgresConfigOption {
	return func(pc *PostgresConfig) {
		pc.Password = password
	}
}

func WithPostgresDBHost(host string) PostgresConfigOption {
	return func(pc *PostgresConfig) {
		pc.Host = host
	}
}

func WithPostgresDBPort(port string) PostgresConfigOption {
	return func(pc *PostgresConfig) {
		pc.Port = port
	}
}

func WithPostgresPoolMaxConns(maxConnection int) PostgresConfigOption {
	return func(pc *PostgresConfig) {
		pc.OptionalConfig["pool_max_conns"] = fmt.Sprint(maxConnection)
	}
}

func WithPostgresPoolMaxConnLifetime(maxConnectionLifetime time.Duration) PostgresConfigOption {
	return func(pc *PostgresConfig) {
		pc.OptionalConfig["pool_max_conn_lifetime"] = maxConnectionLifetime.String()
	}
}

func WithPostgresPoolMaxConnIdleTime(maxConnectionIdleTime time.Duration) PostgresConfigOption {
	return func(pc *PostgresConfig) {
		pc.OptionalConfig["pool_max_conn_idle_time"] = maxConnectionIdleTime.String()
	}
}

type Postgres struct {
	*pgxpool.Pool
}

func NewPostgres(configOpts ...PostgresConfigOption) Postgres {
	config := &PostgresConfig{
		OptionalConfig: make(map[string]string),
	}

	for _, apply := range configOpts {
		apply(config)
	}

	// TODO: config validation

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(config.User, config.Password),
		Host:   fmt.Sprintf("%s:%s", config.Host, config.Port),
		Path:   config.DatabaseName,
	}

	q := dsn.Query()
	q.Add("sslmode", "disable")
	for k, v := range config.OptionalConfig {
		q.Add(k, v)
	}

	dsn.RawQuery = q.Encode()

	dbPool, err := pgxpool.New(context.Background(), dsn.String())
	if err != nil {
		panic(err)
	}

	return Postgres{dbPool}
}

func (p Postgres) Get(ctx context.Context, destination any, query string, args ...any) error {
	return pgxscan.Get(ctx, p.Pool, destination, query, args...)
}

func (p Postgres) Select(ctx context.Context, destination any, query string, args ...any) error {
	return pgxscan.Select(ctx, p.Pool, destination, query, args...)
}
