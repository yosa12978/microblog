package app

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"microblog-app/internal/config"
	"microblog-app/internal/data"
	"microblog-app/internal/session"
	"microblog-app/internal/templates"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	templates  fs.FS
	assets     fs.FS
	migrations fs.FS
	conf       config.Config
	logger     *slog.Logger
}

func New(templates, assets, migrations fs.FS) *App {
	return &App{
		templates:  templates,
		assets:     assets,
		migrations: migrations,
		conf:       config.Get(),
		logger:     slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (app *App) Run() error {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	templates.Init(app.templates)
	if _, err := data.InitPostgres(ctx, app.migrations); err != nil {
		return err
	}
	session.InitStore()

	server := newServer(app.conf.Server.Addr, app.logger, app.assets, app.conf)
	errCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()
	app.logger.Info("server started", "addr", server.Addr)
	var err error
	select {
	case err = <-errCh:
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = server.Shutdown(timeout)
	}
	return err
}
