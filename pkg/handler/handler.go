package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/vekshinnikita/golang_vehicles/docs"
	"github.com/vekshinnikita/golang_vehicles/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"*"}
	router.Use(cors.New(corsConfig))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api")
	{
		vehicles := api.Group("/vehicles")
		{
			vehicles.GET("", h.GetAllVehicles)
			vehicles.POST("", h.CreateVehicle)
			vehicles.PATCH("/:id", h.UpdateVehicle)
			vehicles.DELETE("/:id", h.DeleteVehicle)
		}
	}

	return router
}
