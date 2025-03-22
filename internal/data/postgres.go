package data

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"microblog-app/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	db *pgxpool.Pool
)

func migrateDB(migrations fs.FS, databaseURL string) error {
	dir, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}
	migrator, err := migrate.NewWithSourceInstance("iofs", dir, databaseURL)
	if err != nil {
		return err
	}
	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func InitPostgres(ctx context.Context, migrations fs.FS) (*pgxpool.Pool, error) {
	conf := config.Get()
	addr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		conf.Postgres.User,
		conf.Postgres.Password,
		conf.Postgres.Addr,
		conf.Postgres.DB,
		conf.Postgres.SSLMode,
	)
	conn, err := pgxpool.New(ctx, addr)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}
	db = conn
	return conn, migrateDB(migrations, addr)
}

func Postgres() *pgxpool.Pool {
	return db
}
