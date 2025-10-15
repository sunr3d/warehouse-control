package httphandlers

import (
	"github.com/wb-go/wbf/ginext"

	"github.com/sunr3d/warehouse-control/internal/handlers/middleware"
	"github.com/sunr3d/warehouse-control/internal/interfaces/services"
	"github.com/sunr3d/warehouse-control/models"
)

type handler struct {
	authSvc services.AuthService
	invSvc  services.InventoryService
}

func New(authSvc services.AuthService, invSvc services.InventoryService) *handler {
	return &handler{authSvc: authSvc, invSvc: invSvc}
}

func (h *handler) RegisterHandlers() *ginext.Engine {
	router := ginext.New("")
	router.Use(ginext.Logger(), ginext.Recovery())

	// API
	// Доступны без авторизации
	router.POST("/login", h.login)
	router.GET("/", func(c *ginext.Context) {
		c.File("./web/index.html")
	})
	router.Static("/web", "./web")

	// Доступны только для авторизованных пользователей
	protected := router.Group("/items")
	protected.Use(middleware.AuthMiddleware(h.authSvc))

	protected.GET("", middleware.RBACMiddleware(
		models.RoleAdmin,
		models.RoleManager,
		models.RoleViewer,
	), h.getItems)

	protected.GET("/:id/history", middleware.RBACMiddleware(
		models.RoleAdmin,
		models.RoleManager,
	), h.getItemHistory)

	protected.POST("", middleware.RBACMiddleware(
		models.RoleAdmin,
		models.RoleManager,
	), h.createItem)

	protected.PUT("/:id", middleware.RBACMiddleware(
		models.RoleAdmin,
		models.RoleManager,
	), h.updateItem)

	protected.DELETE("/:id", middleware.RBACMiddleware(
		models.RoleAdmin,
	), h.deleteItem)

	return router
}
