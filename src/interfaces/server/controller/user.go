// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package controller
package controller

import "github.com/labstack/echo/v4"

type UserInterface interface {
	UserRegister(ctx echo.Context) error
}
