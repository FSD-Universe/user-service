// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package repository
package repository

import (
	"time"

	"gorm.io/gorm"
	"half-nothing.cn/service-core/database"
	"half-nothing.cn/service-core/interfaces/database/entity"
	"half-nothing.cn/service-core/interfaces/logger"
)

type RoleRepository struct {
	*database.BaseRepository[*entity.Role]
	pageReq database.PageableInterface[*entity.Role]
}

func NewRoleRepository(
	lg logger.Interface,
	db *gorm.DB,
	queryTimeout time.Duration,
) *RoleRepository {
	return &RoleRepository{
		BaseRepository: database.NewBaseRepository[*entity.Role](lg, "role-repository", db, queryTimeout),
		pageReq:        database.NewPageRequest[*entity.Role](db),
	}
}

func (repo *RoleRepository) SetPermission(roleId uint, permission uint64) error {
	return repo.QueryWithTransaction(func(tx *gorm.DB) error {
		return tx.Model(&entity.Role{}).Where("id = ?", roleId).Update("permission", permission).Error
	})
}

func (repo *RoleRepository) GrantUser(roleId uint, userId uint) error {
	return repo.QueryWithTransaction(func(tx *gorm.DB) error {
		return tx.Create(&entity.UserRole{UserId: userId, RoleId: roleId}).Error
	})
}

func (repo *RoleRepository) RevokeUser(roleId uint, userId uint) error {
	return repo.QueryWithTransaction(func(tx *gorm.DB) error {
		return tx.Delete(&entity.UserRole{}, "user_id = ? AND role_id = ?", userId, roleId).Error
	})
}

func (repo *RoleRepository) GetPages(pageNum int, pageSize int, search string) (roles []*entity.Role, total int64, err error) {
	roles = make([]*entity.Role, 0, pageSize)
	var queryFunc func(tx *gorm.DB) *gorm.DB
	if search != "" {
		queryFunc = func(tx *gorm.DB) *gorm.DB {
			return tx.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
		}
	} else {
		queryFunc = nil
	}
	total, err = repo.QueryWithPagination(repo.pageReq, database.NewPage[*entity.Role](pageNum, pageSize, &roles, &entity.Role{}, queryFunc))
	return
}

func (repo *RoleRepository) GetRoleUsers(roleId uint) (users []*entity.User, err error) {
	users = make([]*entity.User, 0)
	err = repo.QueryWithTransaction(func(tx *gorm.DB) error {
		return tx.Model(&entity.UserRole{}).
			Where("role_id = ?", roleId).
			Joins("User").
			Preload("User.CurrentAvatar").
			Find(&users).
			Error
	})
	return
}

func (repo *RoleRepository) DeleteRole(roleId uint) error {
	return repo.QueryWithTransaction(func(tx *gorm.DB) error {
		err := tx.Delete(&entity.UserRole{}, "role_id = ?", roleId).Error
		if err != nil {
			return err
		}
		return tx.Delete(&entity.Role{ID: roleId}).Error
	})
}
