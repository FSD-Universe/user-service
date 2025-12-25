// Package repository
package repository

import (
	"half-nothing.cn/service-core/interfaces/database/entity"
	"half-nothing.cn/service-core/interfaces/database/repository"
	"half-nothing.cn/service-core/utils"
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
	return userRepo.GetByCid(uint(id))
}

func (id StringUserId) GetUser(userRepo UserInterface) (*entity.User, error) {
	return userRepo.GetByUsernameOrEmail(string(id))
}

type UserInterface interface {
	repository.Base[*entity.User]
	GetByCid(id uint) (*entity.User, error)
	GetByUsernameOrEmail(usernameOrEmail string) (*entity.User, error)
	CheckCidUsernameAndEmail(cid uint, username string, email string) (bool, error)
	GetPages(pageNum int, pageSize int, search string) ([]*entity.User, int64, error)
	GrantRole(userId uint, roleIds []uint) error
	RevokeRole(userId uint, roleIds []uint) error
	GetByIds(userIds []uint) ([]*entity.User, error)
}
