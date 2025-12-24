// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package service
package service

import (
	DTO "user-service/src/interfaces/server/dto"

	"half-nothing.cn/service-core/interfaces/http/dto"
)

type PermissionInterface interface {
	EditUserPermission(data *DTO.EditUserPermission) *dto.ApiResponse[bool]
	EditRolePermission(data *DTO.EditRolePermission) *dto.ApiResponse[bool]
	GrantUserRole(data *DTO.GrantUserRole) *dto.ApiResponse[bool]
	RevokeUserRole(data *DTO.RevokeUserRole) *dto.ApiResponse[bool]
	GrantRoleUser(data *DTO.GrantRoleUser) *dto.ApiResponse[bool]
	RevokeRoleUser(data *DTO.RevokeRoleUser) *dto.ApiResponse[bool]
}
