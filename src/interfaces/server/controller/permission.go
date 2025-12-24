// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package controller
package controller

import "github.com/labstack/echo/v4"

type PermissionInterface interface {
	EditUserPermission(ctx echo.Context) error
	EditRolePermission(ctx echo.Context) error
	GrantUserRole(ctx echo.Context) error
	RevokeUserRole(ctx echo.Context) error
	GrantRoleUser(ctx echo.Context) error
	RevokeRoleUser(ctx echo.Context) error
}
