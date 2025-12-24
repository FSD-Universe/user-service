// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package repository
package repository

import (
	"half-nothing.cn/service-core/interfaces/database/entity"
	"half-nothing.cn/service-core/interfaces/database/repository"
)

type RoleInterface interface {
	repository.Base[*entity.Role]
	GetPages(pageNum int, pageSize int, search string) ([]*entity.Role, int64, error)
	SetPermission(roleId uint, permission uint64) error
	GetRoleUsers(roleId uint) ([]*entity.User, error)
	GrantUser(roleId uint, userId uint) error
	RevokeUser(roleId uint, userId uint) error
	DeleteRole(roleId uint) error
}
