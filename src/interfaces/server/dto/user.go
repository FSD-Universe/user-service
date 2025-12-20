// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package dto
package dto

import (
	"fmt"
	"time"

	"half-nothing.cn/service-core/interfaces/database/entity"
	"half-nothing.cn/service-core/interfaces/http/dto"
)

type BaseUserInfo struct {
	Id            uint      `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Cid           uint      `json:"cid"`
	AvatarUrl     string    `json:"avatar_url"`
	QQ            string    `json:"qq"`
	Rating        int       `json:"rating"`
	Permission    int       `json:"permission"`
	RegisterTime  time.Time `json:"register_time"`
	LastLoginTime time.Time `json:"last_login_time"`
	LastLoginIp   string    `json:"last_login_ip"`
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
	b.Permission = int(user.Permission)
	b.RegisterTime = user.CreatedAt
	if user.LastLoginTime.Valid {
		b.RegisterTime = user.LastLoginTime.Time
	}
	if user.LastLoginIP != nil {
		b.LastLoginIp = *user.LastLoginIP
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

type UserRegister struct {
	dto.HttpContent
	Username string `json:"username" valid:"required,max=64,regex=^[A-Za-z_-][\\w-]*$"`
	Email    string `json:"email" valid:"required,max=128,regex=^[\\w-]+@[\\w-]+(\\.[\\w-]+)+$"`
	Password string `json:"password" valid:"required"`
	Code     string `json:"code" valid:"required,length=6"`
	Cid      int    `json:"cid" valid:"required,min=0"`
}
