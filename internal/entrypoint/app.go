package entrypoint

import (
	"context"
	"fmt"

	"github.com/sunr3d/warehouse-control/internal/config"
	httphandlers "github.com/sunr3d/warehouse-control/internal/handlers"
	"github.com/sunr3d/warehouse-control/internal/infra/postgres"
	"github.com/sunr3d/warehouse-control/internal/interfaces/infra"
	"github.com/sunr3d/warehouse-control/internal/server"
)

func RunApp(ctx context.Context, cfg *config.Config) error {
	// Инфраслой (Infrastructure layer)
	repo, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf("postgres.New: %w", err)
	}
	defer func(db infra.Database) {
		if closer, ok := db.(io.Closer); ok {
			_ = closer.Close()
		}
	}(repo)

	// Сервисный слой (Application / Use Cases layer)
	// TODO: Init svc

	// Слой представления (Presentation layer)
	h := httphandlers.New(svc)
	engine := h.RegisterHandlers()

	// Сервер
	srv := server.New(":"+cfg.HTTPPort, engine)

	return srv.Run(ctx)
}
