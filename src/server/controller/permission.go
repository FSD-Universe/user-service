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

type PermissionController struct {
	logger  logger.Interface
	service service.PermissionInterface
}

func NewPermissionController(
	lg logger.Interface,
	service service.PermissionInterface,
) *PermissionController {
	return &PermissionController{
		logger:  logger.NewLoggerAdapter(lg, "permission-controller"),
		service: service,
	}
}

func (controller *PermissionController) EditUserPermission(ctx echo.Context) error {
	data := &DTO.EditUserPermission{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("EditUserPermission handle fail, bind argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	if data.Data == nil || len(data.Data) == 0 {
		controller.logger.Errorf("EditUserPermission handle fail, argument Data is empty")
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("EditUserPermission handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("EditUserPermission handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	if err := jwt.SetJwtContent(data, ctx); err != nil {
		controller.logger.Errorf("EditUserPermission handle fail, set jwt content err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrUnknownJwtError)
	}
	controller.logger.Debugf("EditUserPermission with argument %#v", data)
	return controller.service.EditUserPermission(data).Response(ctx)
}

func (controller *PermissionController) EditRolePermission(ctx echo.Context) error {
	data := &DTO.EditRolePermission{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("EditRolePermission handle fail, bind argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	if data.Data == nil || len(data.Data) == 0 {
		controller.logger.Errorf("EditRolePermission handle fail, argument Data is empty")
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("EditRolePermission handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("EditRolePermission handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	if err := jwt.SetJwtContent(data, ctx); err != nil {
		controller.logger.Errorf("EditRolePermission handle fail, set jwt content err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrUnknownJwtError)
	}
	controller.logger.Debugf("EditRolePermission with argument %#v", data)
	return controller.service.EditRolePermission(data).Response(ctx)
}

func (controller *PermissionController) GrantUserRole(ctx echo.Context) error {
	data := &DTO.GrantUserRole{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("GrantUserRole handle fail, bind argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("GrantUserRole handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("GrantUserRole handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	if err := jwt.SetJwtContent(data, ctx); err != nil {
		controller.logger.Errorf("GrantUserRole handle fail, set jwt content err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrUnknownJwtError)
	}
	controller.logger.Debugf("GrantUserRole with argument %#v", data)
	return controller.service.GrantUserRole(data).Response(ctx)
}

func (controller *PermissionController) RevokeUserRole(ctx echo.Context) error {
	data := &DTO.RevokeUserRole{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("RevokeUserRole handle fail, bind argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("RevokeUserRole handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("RevokeUserRole handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	if err := jwt.SetJwtContent(data, ctx); err != nil {
		controller.logger.Errorf("RevokeUserRole handle fail, set jwt content err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrUnknownJwtError)
	}
	controller.logger.Debugf("RevokeUserRole with argument %#v", data)
	return controller.service.RevokeUserRole(data).Response(ctx)
}

func (controller *PermissionController) GrantRoleUser(ctx echo.Context) error {
	data := &DTO.GrantRoleUser{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("GrantRoleUser handle fail, bind argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("GrantRoleUser handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("GrantRoleUser handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	if err := jwt.SetJwtContent(data, ctx); err != nil {
		controller.logger.Errorf("GrantRoleUser handle fail, set jwt content err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrUnknownJwtError)
	}
	controller.logger.Debugf("GrantRoleUser with argument %#v", data)
	return controller.service.GrantRoleUser(data).Response(ctx)
}

func (controller *PermissionController) RevokeRoleUser(ctx echo.Context) error {
	data := &DTO.RevokeRoleUser{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("RevokeRoleUser handle fail, bind argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("RevokeRoleUser handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("RevokeRoleUser handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	if err := jwt.SetJwtContent(data, ctx); err != nil {
		controller.logger.Errorf("RevokeRoleUser handle fail, set jwt content err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrUnknownJwtError)
	}
	controller.logger.Debugf("RevokeRoleUser with argument %#v", data)
	return controller.service.RevokeRoleUser(data).Response(ctx)
}
