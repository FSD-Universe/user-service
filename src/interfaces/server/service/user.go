// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package service
package service

import (
	DTO "user-service/src/interfaces/server/dto"

	"half-nothing.cn/service-core/interfaces/http/dto"
)

type UserInterface interface {
	Register(user *DTO.UserRegister) *dto.ApiResponse[bool]
}
