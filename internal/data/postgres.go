package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"microblog-app/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

func migrateDB(migrations fs.FS, conn *sql.DB) error {
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return err
	}
	dir, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}
	migrator, err := migrate.NewWithInstance("iofs", dir, "postgres", driver)
	if err != nil {
		return err
	}
	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func InitPostgres(ctx context.Context, migrations fs.FS) (*sql.DB, error) {
	conf := config.Get()
	addr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		conf.Postgres.User,
		conf.Postgres.Password,
		conf.Postgres.Addr,
		conf.Postgres.DB,
		conf.Postgres.SSLMode,
	)
	conn, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}
	if err := conn.PingContext(ctx); err != nil {
		return nil, err
	}
	db = conn
	return conn, migrateDB(migrations, conn)
}

func Postgres() *sql.DB {
	return db
}
