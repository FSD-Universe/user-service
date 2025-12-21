// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package server
package server

import (
	"io"
	"user-service/src/interfaces/content"
	"user-service/src/server/controller"
	"user-service/src/server/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"half-nothing.cn/service-core/http"
	"half-nothing.cn/service-core/interfaces/logger"
)

func StartServer(content *content.ApplicationContent) {
	c := content.ConfigManager().GetConfig()
	lg := logger.NewLoggerAdapter(content.Logger(), "http-server")

	lg.Info("Http server initializing...")
	e := echo.New()
	e.Logger.SetOutput(io.Discard)
	e.Logger.SetLevel(log.OFF)

	http.SetEchoConfig(lg, e, c.ServerConfig.HttpServerConfig, nil)
	jwtMidware, _, requireRefresh := http.GetJWTMiddleware(content.ClaimFactory())
	if c.TelemetryConfig.HttpServerTrace {
		http.SetTelemetry(e, c.TelemetryConfig)
	}

	authController := controller.NewAuthController(
		content.Logger(),
		service.NewAuthService(
			content.Logger(),
			content.UserRepo(),
			content.ClaimFactory(),
		),
	)

	userController := controller.NewUserController(
		content.Logger(),
		service.NewUserService(
			content.Logger(),
			content.UserRepo(),
			content.EmailClient(),
			content.AuditLogClient(),
		),
	)

	apiGroup := e.Group("/api/v1")
	userGroup := apiGroup.Group("/users")
	userGroup.POST("/token", authController.UserLogin)
	userGroup.POST("/token/fsd", authController.UserFsdLogin)
	userGroup.GET("/token", authController.RefreshToken, jwtMidware, requireRefresh)

	userGroup.POST("", userController.UserRegister)

	http.SetUnmatchedRoute(e)
	http.SetCleaner(content.Cleaner(), e)

	http.Serve(lg, e, c.ServerConfig.HttpServerConfig)
}
