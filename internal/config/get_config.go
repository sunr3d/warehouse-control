package config

import (
	"fmt"

	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/zlog"
)

func GetConfig() (*Config, error) {
	cfg := config.New()
	if err := cfg.Load("config.yml", ".env", ""); err != nil {
		zlog.Logger.Warn().Msgf("config.Load(): %v. Продолжаем с дефолтными значениями...", err)
	}

	cfg.SetDefault("HTTP_PORT", "8080")
	cfg.SetDefault("LOG_LEVEL", "info")
	cfg.SetDefault("JWT_SECRET", "oXzXAjJU0xOuUBtLSfd+zuSU5Xyf1T9wad/IiTrhPT3tri5+cXz1tK7Nizzp7a5lERKKkprGAuZmKYCicbNyQQ") // FIXME: NOT FOR PRODUCTION, JUST FOR DEMO

	cfg.SetDefault("DB_DSN", "postgres://warehouse_control_user:warehouse_control_password@postgres:5432/warehouse_control_db?sslmode=disable")
	cfg.SetDefault("DB_MAX_OPEN_CONNS", 10)
	cfg.SetDefault("DB_MAX_IDLE_CONNS", 2)

	var c Config
	if err := cfg.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("cfg.Unmarshal: %w", err)
	}

	return &c, nil
}
