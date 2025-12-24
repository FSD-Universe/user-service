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

func (controller *UserController) Register(ctx echo.Context) error {
	data := &DTO.UserRegister{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("Register handle fail, parse argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("Register handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("Register handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	controller.logger.Debugf("Register with argument %#v", data)
	return controller.service.Register(data).Response(ctx)
}

func (controller *UserController) CheckAvailability(ctx echo.Context) error {
	data := &DTO.UserCheckAvailability{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("CheckAvailability handle fail, parse argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("CheckAvailability handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("CheckAvailability handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	if data.Cid <= 0 && data.Email == "" && data.Username == "" {
		controller.logger.Errorf("CheckAvailability handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	controller.logger.Debugf("CheckAvailability with argument %#v", data)
	return controller.service.CheckAvailability(data).Response(ctx)
}

func (controller *UserController) ResetPassword(ctx echo.Context) error {
	data := &DTO.UserResetPassword{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("ResetPassword handle fail, parse argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("ResetPassword handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("ResetPassword handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	controller.logger.Debugf("ResetPassword with argument %#v", data)
	return controller.service.ResetPassword(data).Response(ctx)
}

func (controller *UserController) GetPages(ctx echo.Context) error {
	data := &DTO.GetUserPage{}
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

func (controller *UserController) GetSelfData(ctx echo.Context) error {
	data := &DTO.GetCurrentUserData{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("GetSelfData handle fail, parse argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("GetSelfData handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("GetSelfData handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	if err := jwt.SetJwtContent(data, ctx); err != nil {
		controller.logger.Errorf("GetSelfData handle fail, set jwt content err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrUnknownJwtError)
	}
	controller.logger.Debugf("GetSelfData with argument %#v", data)
	return controller.service.GetSelfData(data).Response(ctx)
}

func (controller *UserController) GetData(ctx echo.Context) error {
	data := &DTO.GetUserData{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("GetData handle fail, parse argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("GetData handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("GetData handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	if err := jwt.SetJwtContent(data, ctx); err != nil {
		controller.logger.Errorf("GetData handle fail, set jwt content err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrUnknownJwtError)
	}
	controller.logger.Debugf("GetData with argument %#v", data)
	return controller.service.GetData(data).Response(ctx)
}

func (controller *UserController) UpdateSelfData(ctx echo.Context) error {
	data := &DTO.UpdateCurrentUserData{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("UpdateSelfData handle fail, parse argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	if data.Email == "" && data.Username == "" && data.QQ == "" && data.ImageId == nil {
		controller.logger.Errorf("UpdateSelfData handle fail, nothing need to update")
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	if data.Email != "" && data.EmailCode == "" {
		controller.logger.Errorf("UpdateSelfData handle fail, no email code provided")
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("UpdateSelfData handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("UpdateSelfData handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	if err := jwt.SetJwtContent(data, ctx); err != nil {
		controller.logger.Errorf("UpdateSelfData handle fail, set jwt content err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrUnknownJwtError)
	}
	controller.logger.Debugf("UpdateSelfData with argument %#v", data)
	return controller.service.UpdateSelfData(data).Response(ctx)
}

func (controller *UserController) UpdateData(ctx echo.Context) error {
	data := &DTO.UpdateUserData{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("UpdateData handle fail, parse argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	if data.Email == "" && data.Username == "" && data.QQ == "" && data.Password == "" {
		controller.logger.Errorf("UpdateSelfData handle fail, nothing need to update")
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("UpdateData handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("UpdateData handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, res)
	}
	dto.SetHttpContent(data, ctx)
	if err := jwt.SetJwtContent(data, ctx); err != nil {
		controller.logger.Errorf("UpdateData handle fail, set jwt content err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrUnknownJwtError)
	}
	controller.logger.Debugf("UpdateData with argument %#v", data)
	return controller.service.UpdateData(data).Response(ctx)
}

var ErrSamePassword = dto.NewApiStatus("SAME_PASSWORD", "原密码和新密码不能相同", dto.HttpCodeBadRequest)

func (controller *UserController) UpdatePassword(ctx echo.Context) error {
	data := &DTO.UpdateUserPassword{}
	if err := ctx.Bind(data); err != nil {
		controller.logger.Errorf("UpdatePassword handle fail, parse argument fail, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	if data.NewPassword == data.OldPassword {
		controller.logger.Errorf("UpdatePassword handle fail, new_password cannot be the same as old_password")
		return dto.ErrorResponse(ctx, ErrSamePassword)
	}
	res, err := dto.ValidStruct(data)
	if err != nil {
		controller.logger.Errorf("UpdatePassword handle fail, validate err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrServerError)
	}
	if res != nil {
		controller.logger.Errorf("UpdatePassword handle fail, validate argument fail, %v", res)
		return dto.ErrorResponse(ctx, dto.ErrErrorParam)
	}
	dto.SetHttpContent(data, ctx)
	if err := jwt.SetJwtContent(data, ctx); err != nil {
		controller.logger.Errorf("UpdatePassword handle fail, set jwt content err, %v", err)
		return dto.ErrorResponse(ctx, dto.ErrUnknownJwtError)
	}
	controller.logger.Infof("UpdatePassword with argument %#v", data)
	return controller.service.UpdatePassword(data).Response(ctx)
}
