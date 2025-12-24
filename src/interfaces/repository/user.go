// Package repository
package repository

import (
	"errors"

	"half-nothing.cn/service-core/interfaces/database/entity"
	"half-nothing.cn/service-core/interfaces/database/repository"
	"half-nothing.cn/service-core/utils"
)

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user does not exist")
	// ErrIdentifierTaken 三元组一致性检查失败
	ErrIdentifierTaken = errors.New("user identifiers have been used")
	// ErrIdentifierCheck 三元组一致性检查异常
	ErrIdentifierCheck = errors.New("identifier check error")
	// ErrPasswordEncode 密码编码错误
	ErrPasswordEncode = errors.New("password encode error")
	// ErrOldPassword 原密码错误
	ErrOldPassword = errors.New("old password error")
)

type UserId interface {
	GetUser(userRepo UserInterface) (*entity.User, error)
}

func GetUserId(userId string) UserId {
	id := utils.StrToInt(userId, -1)
	if id == -1 {
		return StringUserId(userId)
	}
	return IntUserId(id)
}

type IntUserId uint
type StringUserId string

func (id IntUserId) GetUser(userRepo UserInterface) (*entity.User, error) {
	return userRepo.GetByIdOrCid(uint(id))
}

func (id StringUserId) GetUser(userRepo UserInterface) (*entity.User, error) {
	return userRepo.GetByUsernameOrEmail(string(id))
}

type UserInterface interface {
	repository.Base[*entity.User]
	GetByIdOrCid(id uint) (*entity.User, error)
	GetByUsernameOrEmail(usernameOrEmail string) (*entity.User, error)
	CheckCidUsernameAndEmail(cid uint, username string, email string) (bool, error)
	GetPages(pageNum int, pageSize int, search string) ([]*entity.User, int64, error)
}
