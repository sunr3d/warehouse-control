package httphandlers

import (
	"github.com/wb-go/wbf/ginext"

	"github.com/sunr3d/warehouse-control/internal/interfaces/services"
)

type handler struct {
	svc services.WarehouseControlService
}

func New(svc services.WarehouseControlService) *handler {
	return &handler{svc: svc}
}

func (h *handler) RegisterHandlers() *ginext.Engine {
	router := ginext.New("")
	router.Use(ginext.Logger(), ginext.Recovery())

	// API
	// TODO: Add handlers

	// Web-UI
	router.Static("/web", "./web")
	router.GET("/", func(c *ginext.Context) {
		c.File("./web/index.html")
	})

	return router
}
