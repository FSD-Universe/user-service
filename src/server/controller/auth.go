// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package controller
package controller

import (
	DTO "user-service/src/interfaces/server/dto"
	"user-service/src/interfaces/server/service"

	"github.com/labstack/echo/v4"
	"half-nothing.cn/service-core/interfaces/http/dto"
	"half-nothing.cn/service-core/interfaces/http/jwt"
	"half-nothing.cn/service-core/interfaces/logger"
)

type AuthController struct {
	logger      logger.Interface
	userService service.AuthInterface
}

func NewAuthController(
	lg logger.Interface,
	userService service.AuthInterface,
) *AuthController {
	return &AuthController{
		logger:      logger.NewLoggerAdapter(lg, "user-controller"),
		userService: userService,
	}
}

func (controller *AuthController) UserLogin(ctx echo.Context) error {
	data := &DTO.UserLogin{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("UserLogin handle fail, parse argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("UserLogin handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("UserLogin handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	controller.logger.Debugf("UserLogin with argument %#v", data)
	return controller.userService.Login(data).Response(ctx)
}

func (controller *AuthController) UserFsdLogin(ctx echo.Context) error {
	data := &DTO.UserFsdLogin{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("UserFsdLogin handle fail, parse argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("UserFsdLogin handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("UserFsdLogin handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	controller.logger.Debugf("UserFsdLogin with argument %#v", data)
	return dto.JsonResponse(ctx, dto.HttpCodeOk.Code(), controller.userService.FsdLogin(data))
}

func (controller *AuthController) RefreshToken(ctx echo.Context) error {
	data := &DTO.RefreshToken{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("RefreshToken handle fail, parse argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	dto.SetHttpContent(data, ctx)
	err := jwt.SetJwtContent(data, ctx)
	if err != nil {
		controller.logger.Errorf("RefreshToken handle fail, set jwt content err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	controller.logger.Debugf("RefreshToken with argument %#v", data)
	return controller.userService.RefreshToken(data).Response(ctx)
}
