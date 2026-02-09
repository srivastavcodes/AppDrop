package main

import (
	"appdrop/internal/data"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type config struct {
	port int
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleTime  time.Duration
		maxIdleConns int
	}
	cors struct {
		allowedOrigins []string
	}
}

type backend struct {
	logger *slog.Logger
	conf   config
	models data.Models
	wg     sync.WaitGroup
}

func (b *backend) serve() error {
	srv := &http.Server{
		ErrorLog:          slog.NewLogLogger(b.logger.Handler(), slog.LevelError),
		Addr:              fmt.Sprintf(":%d", b.conf.port),
		Handler:           b.routes(),
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       10 * time.Second,
		IdleTimeout:       90 * time.Second,
	}
	shutdownErr := make(chan error, 1)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sigv := <-quit
		b.logger.Info("shutting down server", "signal", sigv)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			shutdownErr <- err
		}
		b.logger.Info("completing background tasks", "addr", srv.Addr)
		b.wg.Wait()
		shutdownErr <- nil
	}()
	b.logger.Info("server started", "addr", srv.Addr)

	err := srv.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	if err = <-shutdownErr; err != nil {
		return err
	}
	b.logger.Info("server stopped", "addr", srv.Addr)
	return nil
}
