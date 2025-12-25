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

func (repo *RoleRepository) GetRoleUsers(roleId uint) (userRoles []*entity.UserRole, err error) {
	userRoles = make([]*entity.UserRole, 0)
	err = repo.QueryWithTransaction(func(tx *gorm.DB) error {
		return tx.Model(&entity.UserRole{}).
			Where("role_id = ?", roleId).
			Joins("User").
			Joins("User.CurrentAvatar").
			Find(&userRoles).
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

func (repo *RoleRepository) GetByIds(roleIds []uint) (roles []*entity.Role, err error) {
	roles = make([]*entity.Role, 0, len(roleIds))
	err = repo.QueryWithTransaction(func(tx *gorm.DB) error {
		return tx.Find(&roles, roleIds).Error
	})
	return
}

func (repo *RoleRepository) GrantUser(roleId uint, userIds []uint) error {
	userRoles := make([]*entity.UserRole, len(userIds))
	for i, userId := range userIds {
		userRoles[i] = &entity.UserRole{UserId: userId, RoleId: roleId}
	}
	return repo.QueryWithTransaction(func(tx *gorm.DB) error {
		return tx.Create(userRoles).Error
	})
}

func (repo *RoleRepository) RevokeUser(roleId uint, userIds []uint) error {
	return repo.QueryWithTransaction(func(tx *gorm.DB) error {
		return tx.Delete(&entity.UserRole{}, "role_id = ? AND user_id IN ?", roleId, userIds).Error
	})
}
