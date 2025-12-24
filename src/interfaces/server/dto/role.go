// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package dto
package dto

import (
	"half-nothing.cn/service-core/interfaces/database/entity"
	"half-nothing.cn/service-core/interfaces/http/dto"
	"half-nothing.cn/service-core/interfaces/http/jwt"
)

type BaseRoleInfo struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Permission  uint64 `json:"permission"`
}

func (role *BaseRoleInfo) FromRoleEntity(entity *entity.Role) {
	role.Id = entity.ID
	role.Name = entity.Name
	role.Description = entity.Comment
	role.Permission = entity.Permission
}

type RoleInfo struct {
	BaseRoleInfo
	Users []*BaseUserInfo `json:"users"`
}

func (role *RoleInfo) FromRoleEntity(entity *entity.Role) {
	role.BaseRoleInfo.FromRoleEntity(entity)
}

type CreateRole struct {
	dto.HttpContent
	jwt.Content
	Name        string `json:"name" valid:"required"`
	Description string `json:"description"`
}

type UpdateRole struct {
	dto.HttpContent
	jwt.Content
	Id          uint   `param:"id" valid:"required,min=0;exclude"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DeleteRole struct {
	dto.HttpContent
	jwt.Content
	Id    uint `param:"id" valid:"required,min=0;exclude"`
	Force bool `query:"force"`
}

type GetRoleDetail struct {
	dto.HttpContent
	jwt.Content
	Id uint `param:"id" valid:"required,min=0;exclude"`
}

type GetRolePage struct {
	dto.HttpContent
	jwt.Content
	PageNum  int    `query:"page_num" valid:"required,min=0;exclude"`
	PageSize int    `query:"page_size" valid:"required,min=0;exclude"`
	Search   string `query:"search"`
}

type GetRolePageResponse struct {
	Data     []*BaseRoleInfo `json:"page_data"`
	Total    int             `json:"total"`
	PageNum  int             `json:"page_num"`
	PageSize int             `json:"page_size"`
}
