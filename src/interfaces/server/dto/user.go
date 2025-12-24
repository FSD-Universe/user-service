// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package dto
package dto

import (
	"fmt"
	"time"

	"half-nothing.cn/service-core/interfaces/database/entity"
	"half-nothing.cn/service-core/interfaces/http/dto"
	"half-nothing.cn/service-core/interfaces/http/jwt"
	"half-nothing.cn/service-core/permission"
	"half-nothing.cn/service-core/utils"
)

type BaseUserInfo struct {
	Id            uint       `json:"id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	Cid           uint       `json:"cid"`
	AvatarUrl     string     `json:"avatar_url"`
	QQ            string     `json:"qq"`
	Rating        int        `json:"rating"`
	Permission    int        `json:"permission"`
	RegisterTime  time.Time  `json:"register_time"`
	LastLoginTime *time.Time `json:"last_login_time"`
	LastLoginIp   *string    `json:"last_login_ip"`
}

func (b *BaseUserInfo) FromUserEntity(user *entity.User) *BaseUserInfo {
	b.Id = user.ID
	b.Username = user.Username
	b.Email = user.Email
	b.Cid = user.Cid
	if user.QQ != nil {
		b.QQ = *user.QQ
	}
	b.Rating = user.Rating
	perm := permission.Permission(user.Permission)
	utils.ForEach(user.Roles, func(index int, role *entity.UserRole) {
		perm.Merge(permission.Permission(role.Role.Permission))
	})
	b.Permission = int(perm)
	b.RegisterTime = user.CreatedAt
	if user.LastLoginTime.Valid {
		b.LastLoginTime = &user.LastLoginTime.Time
	}
	if user.LastLoginIP != nil {
		b.LastLoginIp = user.LastLoginIP
	}
	if user.CurrentAvatar != nil {
		b.AvatarUrl = user.CurrentAvatar.Url
	} else if b.QQ != "" {
		b.AvatarUrl = fmt.Sprintf("https://q2.qlogo.cn/headimg_dl?dst_uin=%s&spec=100", b.QQ)
	} else {
		b.AvatarUrl = ""
	}
	return b
}

type UserInfo struct {
	BaseUserInfo
	Banned     bool            `json:"banned"`
	BannedTime *time.Time      `json:"banned_time"`
	Roles      []*BaseRoleInfo `json:"roles"`
}

func (u *UserInfo) FromUserEntity(user *entity.User) *UserInfo {
	u.BaseUserInfo.FromUserEntity(user)
	u.Banned = user.Banned
	if user.BannedUntil.Valid {
		u.BannedTime = &user.BannedUntil.Time
	}
	u.Roles = make([]*BaseRoleInfo, len(user.Roles))
	utils.ForEach(user.Roles, func(index int, role *entity.UserRole) {
		u.Roles[index] = &BaseRoleInfo{}
		u.Roles[index].FromRoleEntity(role.Role)
	})
	return u
}

type UserRegister struct {
	dto.HttpContent
	Username string `json:"username" valid:"required,max=64,regex=^[A-Za-z_-][\\w-]*$"`
	Email    string `json:"email" valid:"required,max=128,regex=^[\\w-]+@[\\w-]+(\\.[\\w-]+)+$"`
	Password string `json:"password" valid:"required"`
	Code     string `json:"code" valid:"required,length=6"`
	Cid      int    `json:"cid" valid:"required,min=0;exclude"`
}

type UserCheckAvailability struct {
	dto.HttpContent
	Username string `query:"username" valid:"max=64,regex=^[A-Za-z_-][\\w-]*$"`
	Email    string `query:"email" valid:"max=128,regex=^[\\w-]+@[\\w-]+(\\.[\\w-]+)+$"`
	Cid      int    `query:"cid"`
}

type UserResetPassword struct {
	dto.HttpContent

	Email    string `json:"email" valid:"required,max=128,regex=^[\\w-]+@[\\w-]+(\\.[\\w-]+)+$"`
	Code     string `json:"code" valid:"required,length=6"`
	Password string `json:"password" valid:"required"`
}

type GetUserPage struct {
	dto.HttpContent
	jwt.Content
	PageNum  int    `query:"page_num" valid:"required,min=0;exclude"`
	PageSize int    `query:"page_size" valid:"required,min=0;exclude"`
	Search   string `query:"search"`
}

type GetUserPageResponse struct {
	Data     []*UserInfo `json:"page_data"`
	Total    int         `json:"total"`
	PageNum  int         `json:"page_num"`
	PageSize int         `json:"page_size"`
}

type GetCurrentUserData struct {
	dto.HttpContent
	jwt.Content
}

type GetUserData struct {
	dto.HttpContent
	jwt.Content
	Id uint `param:"id" valid:"required,min=0;exclude"`
}

type UpdateCurrentUserData struct {
	dto.HttpContent
	jwt.Content
	Username  string `json:"username" valid:"max=64,regex=^[A-Za-z_-][\\w-]*$"`
	Email     string `json:"email" valid:"max=128,regex=^[\\w-]+@[\\w-]+(\\.[\\w-]+)+$"`
	EmailCode string `json:"email_code" valid:"length=6"`
	QQ        string `json:"qq" valid:"max=16,regex=^[1-9][0-9]*$"`
	ImageId   *uint  `json:"image_id"`
}

type UpdateUserPassword struct {
	dto.HttpContent
	jwt.Content
	OldPassword string `json:"old_password" valid:"required"`
	NewPassword string `json:"new_password" valid:"required"`
}

type UpdateUserData struct {
	dto.HttpContent
	jwt.Content
	Id       uint   `param:"id" valid:"required,min=0;exclude"`
	Username string `json:"username" valid:"max=64,regex=^[A-Za-z_-][\\w-]*$"`
	Email    string `json:"email" valid:"max=128,regex=^[\\w-]+@[\\w-]+(\\.[\\w-]+)+$"`
	QQ       string `json:"qq" valid:"max=16,regex=^[1-9][0-9]*$"`
	Password string `json:"password"`
}
