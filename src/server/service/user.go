// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package service
package service

import (
	"user-service/src/interfaces/repository"
	DTO "user-service/src/interfaces/server/dto"

	"half-nothing.cn/service-core/interfaces/http/dto"
	"half-nothing.cn/service-core/interfaces/logger"
)

type UserService struct {
	logger logger.Interface
	repo   repository.UserInterface
}

func NewUserService(
	lg logger.Interface,
	repo repository.UserInterface,
) *UserService {
	return &UserService{
		logger: logger.NewLoggerAdapter(lg, "user-service"),
		repo:   repo,
	}
}

func (u *UserService) Register(user *DTO.UserRegister) *dto.ApiResponse[bool] {
	//TODO implement me
	panic("implement me")
}
