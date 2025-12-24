// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package dto
package dto

import (
	"half-nothing.cn/service-core/interfaces/http/dto"
	"half-nothing.cn/service-core/interfaces/http/jwt"
)

type EditUserPermission struct {
	dto.HttpContent
	jwt.Content
	UserId uint            `param:"id" valid:"required,min=0;exclude"`
	Data   map[string]bool `json:"data"`
}

type EditRolePermission struct {
	dto.HttpContent
	jwt.Content
	RoleId uint            `param:"id" valid:"required,min=0;exclude"`
	Data   map[string]bool `json:"data"`
}

type GrantRoleUser struct {
	dto.HttpContent
	jwt.Content
	RoleId  uint   `param:"id" valid:"required,min=0;exclude"`
	UserIds []uint `json:"ids" valid:"required,min=0;exclude"`
}

type RevokeRoleUser struct {
	dto.HttpContent
	jwt.Content
	RoleId  uint   `param:"id" valid:"required,min=0;exclude"`
	UserIds []uint `json:"ids" valid:"required,min=0;exclude"`
}

type GrantUserRole struct {
	dto.HttpContent
	jwt.Content
	UserId  uint   `param:"id" valid:"required,min=0;exclude"`
	RoleIds []uint `json:"ids" valid:"required,min=0;exclude"`
}

type RevokeUserRole struct {
	dto.HttpContent
	jwt.Content
	UserId  uint   `param:"id" valid:"required,min=0;exclude"`
	RoleIds []uint `json:"ids" valid:"required,min=0;exclude"`
}
