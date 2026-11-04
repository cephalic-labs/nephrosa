package database

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/cephalic-labs/nephrosa/backend/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool = pgxpool.Pool

func Open(ctx context.Context, cfg config.Database, simpleProtocol bool) (*Pool, error) {
	passwordBytes, err := os.ReadFile(cfg.PasswordFile)
	if err != nil {
		return nil, fmt.Errorf("read database credential file: %w", err)
	}
	password := strings.TrimRight(string(passwordBytes), "\r\n")
	if password == "" {
		return nil, errors.New("database credential file is empty")
	}

	connectionURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, password),
		Host:   net.JoinHostPort(cfg.Host, fmt.Sprint(cfg.Port)),
		Path:   "/" + cfg.Name,
		RawQuery: url.Values{
			"sslmode": {cfg.SSLMode},
		}.Encode(),
	}
	poolConfig, err := pgxpool.ParseConfig(connectionURL.String())
	if err != nil {
		return nil, errors.New("invalid database connection settings")
	}
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 0
	poolConfig.ConnConfig.ConnectTimeout = 5 * time.Second
	if simpleProtocol {
		poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create database connection pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("connect to PostgreSQL: %w", err)
	}
	return pool, nil
}
