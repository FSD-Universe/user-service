// Package controller
package controller

import (
	"github.com/labstack/echo/v4"
)

type AuthInterface interface {
	UserLogin(ctx echo.Context) error
	UserFsdLogin(ctx echo.Context) error
	UserRegister(ctx echo.Context) error
	RefreshToken(ctx echo.Context) error
}
