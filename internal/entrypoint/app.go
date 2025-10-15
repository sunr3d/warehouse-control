package entrypoint

import (
	"context"
	"fmt"
	"io"

	"github.com/sunr3d/warehouse-control/internal/config"
	httphandlers "github.com/sunr3d/warehouse-control/internal/handlers"
	"github.com/sunr3d/warehouse-control/internal/infra/postgres"
	"github.com/sunr3d/warehouse-control/internal/interfaces/infra"
	"github.com/sunr3d/warehouse-control/internal/server"
	"github.com/sunr3d/warehouse-control/internal/services/authsvc"
	"github.com/sunr3d/warehouse-control/internal/services/inventorysvc"
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
	authSvc := authsvc.New(repo, cfg.JWTSecret)
	invSvc := inventorysvc.New(repo)

	// Слой представления (Presentation layer)
	h := httphandlers.New(authSvc, invSvc)
	engine := h.RegisterHandlers()

	// Сервер
	srv := server.New(":"+cfg.HTTPPort, engine)

	return srv.Run(ctx)
}
