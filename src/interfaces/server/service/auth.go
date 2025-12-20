// Package service
package service

import (
	"user-service/src/interfaces/server/dto"

	. "half-nothing.cn/service-core/interfaces/http/dto"
)

var (
	ErrUsernameOrPasswordError = NewApiStatus("LOGIN_FAIL", "用户名或密码错误", HttpCodeBadRequest)
	ErrUserBanned              = NewApiStatus("USER_BANNED", "您已被封禁", HttpCodePermissionDenied)
)

type AuthInterface interface {
	Login(form *dto.UserLogin) *ApiResponse[*dto.UserLoginResponse]
	FsdLogin(form *dto.UserFsdLogin) *dto.UserFsdLoginResponse
	RefreshToken(form *dto.RefreshToken) *ApiResponse[*dto.RefreshTokenResponse]
}
