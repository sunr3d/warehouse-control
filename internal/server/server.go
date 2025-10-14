package server

import (
	"context"
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

const (
	RWTimeout       = 15 * time.Second
	IdleTimeout     = 60 * time.Second
	ShutdownTimeout = 30 * time.Second
)

type server struct {
	engine *ginext.Engine
	addr   string
}

// New - конструктор Server.
func New(addr string, engine *ginext.Engine) *server {
	return &server{
		engine: engine,
		addr:   addr,
	}
}

// Run - запуск HTTP сервера с обработкой graceful shutdown.
func (s *server) Run(ctx context.Context) error {
	httpServer := &http.Server{
		Addr:         s.addr,
		Handler:      s.engine,
		ReadTimeout:  RWTimeout,
		WriteTimeout: RWTimeout,
		IdleTimeout:  IdleTimeout,
	}

	srvErr := make(chan error, 1)
	go func() {
		zlog.Logger.Info().Msgf("HTTP сервер запущен на %s", s.addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			srvErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		zlog.Logger.Info().Msg("Получен сигнал о завершении работы, инициализация graceful shutdown...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			zlog.Logger.Error().Err(err).Msg("Ошибка при завершении работы HTTP сервера")
			return err
		}

		zlog.Logger.Info().Msg("HTTP сервер завершен успешно")
		return nil

	case err := <-srvErr:
		zlog.Logger.Error().Err(err).Msg("HTTP сервер завершен с ошибкой")
		return err
	}
}
