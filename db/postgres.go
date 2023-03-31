package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
	"todoproject/apperror"
)

type Config struct {
	Username     string
	Password     string
	Host         string
	Port         string
	Database     string
	TokenExpires time.Duration
	TokenSecret  string
	TokenMaxAge  int
}

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

// NewClient "postgres://username:password@localhost:5432/database_name"
func NewClient(cfg *Config) (*pgx.Conn, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	connect, errConnect := pgx.Connect(context.Background(), url)

	if errConnect != nil {
		return nil, apperror.NewAppError(fmt.Sprintf("failed to connect database. due to error: %v", errConnect))
	}
	return connect, nil
}
