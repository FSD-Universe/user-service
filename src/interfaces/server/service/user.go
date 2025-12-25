// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package service
package service

import (
	DTO "user-service/src/interfaces/server/dto"

	"half-nothing.cn/service-core/interfaces/http/dto"
)

type UserInterface interface {
	Register(data *DTO.UserRegister) *dto.ApiResponse[bool]
	CheckAvailability(data *DTO.UserCheckAvailability) *dto.ApiResponse[bool]
	ResetPassword(data *DTO.UserResetPassword) *dto.ApiResponse[bool]
	GetPages(data *DTO.GetUserPage) *dto.ApiResponse[*DTO.GetUserPageResponse]
	GetSelfData(data *DTO.GetCurrentUserData) *dto.ApiResponse[*DTO.UserInfo]
	GetData(data *DTO.GetUserData) *dto.ApiResponse[*DTO.FullUserInfo]
	UpdateSelfData(data *DTO.UpdateCurrentUserData) *dto.ApiResponse[*DTO.UserInfo]
	UpdateData(data *DTO.UpdateUserData) *dto.ApiResponse[bool]
	UpdatePassword(data *DTO.UpdateUserPassword) *dto.ApiResponse[bool]
}
