// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package dto
package dto

import (
	"half-nothing.cn/service-core/interfaces/http/dto"
	"half-nothing.cn/service-core/interfaces/http/jwt"
)

type UserLogin struct {
	dto.HttpContent
	Username string `json:"username" valid:"required,max=128"`
	Password string `json:"password" valid:"required"`
}

type UserLoginResponse struct {
	User         *BaseUserInfo `json:"user"`
	Token        string        `json:"token"`
	ExpiresIn    int           `json:"expires_in"`
	RefreshToken string        `json:"refresh_token"`
}

type UserFsdLogin struct {
	dto.HttpContent
	Cid        string `json:"cid" valid:"required,min=0"`
	Password   string `json:"password" valid:"required"`
	IsSweatbox bool   `json:"is_sweatbox"`
}

type UserFsdLoginResponse struct {
	Success  bool   `json:"success"`
	ErrorMsg string `json:"error_msg"`
	Token    string `json:"token,omitempty"`
}

type RefreshToken struct {
	dto.HttpContent
	jwt.Content
	Force bool `query:"force"`
}

type RefreshTokenResponse = UserLoginResponse
