// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package service
package service

import (
	DTO "user-service/src/interfaces/server/dto"

	"half-nothing.cn/service-core/interfaces/http/dto"
)

type RoleInterface interface {
	GetPages(page *DTO.GetRolePage) *dto.ApiResponse[*DTO.GetRolePageResponse]
	GetById(data *DTO.GetRoleDetail) *dto.ApiResponse[*DTO.RoleInfo]
	Create(role *DTO.CreateRole) *dto.ApiResponse[bool]
	Update(role *DTO.UpdateRole) *dto.ApiResponse[bool]
	Delete(role *DTO.DeleteRole) *dto.ApiResponse[bool]
}
