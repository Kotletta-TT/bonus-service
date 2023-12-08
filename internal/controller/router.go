package controller

import (
	"github.com/Kotletta-TT/bonus-service/config"
	"github.com/Kotletta-TT/bonus-service/internal/controller/middlewares"
	"github.com/Kotletta-TT/bonus-service/internal/controller/routes"
	"github.com/Kotletta-TT/bonus-service/internal/logger"
	"github.com/Kotletta-TT/bonus-service/internal/repository"
	"github.com/gin-gonic/gin"
)

func Router(config *config.Config, repo repository.Repository) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.RedirectTrailingSlash = false
	engine.Use(middlewares.RequestResponseLogging())
	// TODO есть не json-сериализуемые объекты
	// engine.Use(gzip)
	engine.POST("/api/user/register", routes.RegisterHandler(repo, config))
	engine.POST("/api/user/login", routes.LoginHandler(repo, config))
	engine.POST("/api/user/orders", middlewares.Auth(config, repo), routes.OrderSetHandler(repo))
	engine.GET("/api/user/orders", middlewares.Auth(config, repo), routes.OrdersListHandler(repo))
	engine.GET("/api/user/balance", middlewares.Auth(config, repo), routes.GetBalanceHandler(repo))
	engine.POST("/api/user/balance/withdraw", middlewares.Auth(config, repo), routes.RequestWithDrawHandler(repo))
	engine.GET("/api/user/withdrawals", middlewares.Auth(config, repo), routes.GetWithDrawalsHandler(repo))
	logger.Info("Run server", "address:", config.ServAddr)
	return engine
}
