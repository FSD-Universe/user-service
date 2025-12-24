// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package controller
package controller

import "github.com/labstack/echo/v4"

type UserInterface interface {
	Register(ctx echo.Context) error
	CheckAvailability(ctx echo.Context) error
	ResetPassword(ctx echo.Context) error
	GetPages(ctx echo.Context) error
	GetSelfData(ctx echo.Context) error
	GetData(ctx echo.Context) error
	UpdateSelfData(ctx echo.Context) error
	UpdateData(ctx echo.Context) error
	UpdatePassword(ctx echo.Context) error
}
