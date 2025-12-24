// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package controller
package controller

import "github.com/labstack/echo/v4"

type RoleInterface interface {
	GetPages(ctx echo.Context) error
	GetById(ctx echo.Context) error
	Create(ctx echo.Context) error
	Update(ctx echo.Context) error
	Delete(ctx echo.Context) error
}
