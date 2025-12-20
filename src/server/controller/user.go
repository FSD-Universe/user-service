// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package controller
package controller

import (
	DTO "user-service/src/interfaces/server/dto"
	"user-service/src/interfaces/server/service"

	"github.com/labstack/echo/v4"
	"half-nothing.cn/service-core/interfaces/http/dto"
	"half-nothing.cn/service-core/interfaces/logger"
)

type UserController struct {
	logger  logger.Interface
	service service.UserInterface
}

func NewUserController(
	lg logger.Interface,
	service service.UserInterface,
) *UserController {
	return &UserController{
		logger:  logger.NewLoggerAdapter(lg, "user-controller"),
		service: service,
	}
}

func (controller *UserController) UserRegister(ctx echo.Context) error {
	data := &DTO.UserRegister{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("UserRegister handle fail, parse argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("UserRegister handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("UserRegister handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	controller.logger.Debugf("UserRegister with argument %#v", data)
	return controller.service.Register(data).Response(ctx)
}
