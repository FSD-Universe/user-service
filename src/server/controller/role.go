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

type RoleController struct {
	logger  logger.Interface
	service service.RoleInterface
}

func NewRoleController(
	lg logger.Interface,
	service service.RoleInterface,
) *RoleController {
	return &RoleController{
		logger:  logger.NewLoggerAdapter(lg, "role-controller"),
		service: service,
	}
}

func (controller *RoleController) GetPages(ctx echo.Context) error {
	data := &DTO.GetRolePage{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("GetPages handle fail, parse argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("GetPages handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("GetPages handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	if err := jwt.SetJwtContent(data, ctx); err != nil {
		controller.logger.Errorf("GetPages handle fail, set jwt content err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrUnknownJwtError)
	}
	controller.logger.Debugf("GetPages with argument %#v", data)
	return controller.service.GetPages(data).Response(ctx)
}

func (controller *RoleController) GetById(ctx echo.Context) error {
	data := &DTO.GetRoleDetail{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("GetById handle fail, parse argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("GetById handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("GetById handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	if err := jwt.SetJwtContent(data, ctx); err != nil {
		controller.logger.Errorf("GetById handle fail, set jwt content err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrUnknownJwtError)
	}
	controller.logger.Debugf("GetById with argument %#v", data)
	return controller.service.GetById(data).Response(ctx)
}

func (controller *RoleController) Create(ctx echo.Context) error {
	data := &DTO.CreateRole{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("Create handle fail, parse argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("Create handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("Create handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	if err := jwt.SetJwtContent(data, ctx); err != nil {
		controller.logger.Errorf("Create handle fail, set jwt content err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrUnknownJwtError)
	}
	controller.logger.Debugf("Create with argument %#v", data)
	return controller.service.Create(data).Response(ctx)
}

func (controller *RoleController) Update(ctx echo.Context) error {
	data := &DTO.UpdateRole{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("Update handle fail, parse argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	if data.Name == "" && data.Description == "" {
		controller.logger.Error("Update handle fail, nothing to update")
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("Update handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("Update handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	if err := jwt.SetJwtContent(data, ctx); err != nil {
		controller.logger.Errorf("Update handle fail, jwt content err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrUnknownJwtError)
	}
	controller.logger.Debugf("Update with argument %#v", data)
	return controller.service.Update(data).Response(ctx)
}

func (controller *RoleController) Delete(ctx echo.Context) error {
	data := &DTO.DeleteRole{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("Delete handle fail, parse argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("Delete handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("Delete handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	if err := jwt.SetJwtContent(data, ctx); err != nil {
		controller.logger.Errorf("Delete handle fail, set jwt content err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrUnknownJwtError)
	}
	controller.logger.Debugf("Delete with argument %#v", data)
	return controller.service.Delete(data).Response(ctx)
}
