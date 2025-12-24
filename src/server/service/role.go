// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package service
package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
	"user-service/src/interfaces/grpc"
	"user-service/src/interfaces/repository"
	DTO "user-service/src/interfaces/server/dto"

	"gorm.io/gorm"
	"half-nothing.cn/service-core/interfaces/database/entity"
	"half-nothing.cn/service-core/interfaces/http/dto"
	"half-nothing.cn/service-core/interfaces/logger"
	"half-nothing.cn/service-core/permission"
	"half-nothing.cn/service-core/utils"
)

type RoleService struct {
	logger         logger.Interface
	repo           repository.RoleInterface
	auditLogClient grpc.AuditLogClient
}

func NewRoleService(
	lg logger.Interface,
	repo repository.RoleInterface,
	auditLogClient grpc.AuditLogClient,
) *RoleService {
	return &RoleService{
		logger:         logger.NewLoggerAdapter(lg, "role-service"),
		repo:           repo,
		auditLogClient: auditLogClient,
	}
}

func (service *RoleService) GetPages(page *DTO.GetRolePage) *dto.ApiResponse[*DTO.GetRolePageResponse] {
	perm := permission.Permission(page.Permission)
	if !perm.HasPermission(permission.RoleShowList) {
		return dto.NewApiResponse[*DTO.GetRolePageResponse](dto.ErrNoPermission, nil)
	}
	roles, total, err := service.repo.GetPages(page.PageNum, page.PageSize, page.Search)
	if err != nil {
		service.logger.Errorf("error occurred when get role list: %v", err)
		return dto.NewApiResponse[*DTO.GetRolePageResponse](ErrDataBaseError, nil)
	}
	roleInfos := make([]*DTO.BaseRoleInfo, len(roles))
	utils.ForEach(roles, func(index int, role *entity.Role) {
		roleInfos[index] = &DTO.BaseRoleInfo{}
		roleInfos[index].FromRoleEntity(role)
	})
	return dto.NewApiResponse(dto.SuccessHandleRequest, &DTO.GetRolePageResponse{
		Data:     roleInfos,
		Total:    int(total),
		PageNum:  page.PageNum,
		PageSize: page.PageSize,
	})
}

func (service *RoleService) GetById(data *DTO.GetRoleDetail) *dto.ApiResponse[*DTO.RoleInfo] {
	perm := permission.Permission(data.Permission)
	if !perm.HasPermission(permission.RoleShowList) {
		service.logger.Errorf("user %04d no permission to show role list", data.Cid)
		return dto.NewApiResponse[*DTO.RoleInfo](dto.ErrNoPermission, nil)
	}
	role, err := service.repo.GetById(data.Id)
	if err != nil {
		service.logger.Errorf("error occurred when get role by id: %v", err)
		return dto.NewApiResponse[*DTO.RoleInfo](ErrDataBaseError, nil)
	}
	roleInfo := &DTO.RoleInfo{}
	roleInfo.FromRoleEntity(role)
	users, err := service.repo.GetRoleUsers(role.ID)
	roleInfo.Users = make([]*DTO.BaseUserInfo, len(users))
	utils.ForEach(users, func(index int, user *entity.User) {
		roleInfo.Users[index] = &DTO.BaseUserInfo{}
		roleInfo.Users[index].FromUserEntity(user)
	})
	return dto.NewApiResponse(dto.SuccessHandleRequest, roleInfo)
}

//goland:noinspection DuplicatedCode
func (service *RoleService) Create(role *DTO.CreateRole) *dto.ApiResponse[bool] {
	perm := permission.Permission(role.Permission)
	if !perm.HasPermission(permission.RoleCreate) {
		service.logger.Errorf("user %04d no permission to create role", role.Cid)
		return dto.NewApiResponse(dto.ErrNoPermission, false)
	}

	roleEntity := &entity.Role{
		Name:    role.Name,
		Comment: role.Description,
	}
	if err := service.repo.Save(roleEntity); err != nil {
		service.logger.Errorf("error occurred when create role: %v", err)
		return dto.NewApiResponse(ErrDataBaseError, false)
	}

	go func(role *DTO.CreateRole, roleEntity *entity.Role) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := service.auditLogClient.Log(ctx, &grpc.AuditLogRequest{
			Event:     entity.AuditEventRoleCreated.Value,
			Subject:   fmt.Sprintf("%04d", role.Cid),
			Object:    strconv.Itoa(int(roleEntity.ID)),
			Ip:        role.Ip,
			UserAgent: role.UserAgent,
			NewValue:  fmt.Sprintf("%s(%s)", roleEntity.Name, roleEntity.Comment),
		})
		if err != nil {
			service.logger.Errorf("error occurred when create audit log: %v", err)
		}
	}(role, roleEntity)

	return dto.NewApiResponse(dto.SuccessHandleRequest, true)
}

var (
	ErrRoleNotFound = dto.NewApiStatus("ROLE_NOT_FOUND", "角色不存在", dto.HttpCodeNotFound)
)

//goland:noinspection DuplicatedCode
func (service *RoleService) Update(role *DTO.UpdateRole) *dto.ApiResponse[bool] {
	perm := permission.Permission(role.Permission)
	if !perm.HasPermission(permission.RoleEdit) {
		service.logger.Errorf("user %04d no permission to edit role", role.Cid)
		return dto.NewApiResponse(dto.ErrNoPermission, false)
	}

	roleEntity, err := service.repo.GetById(role.Id)
	if err != nil {
		service.logger.Errorf("error occurred when get role by id: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.NewApiResponse(ErrRoleNotFound, false)
		}
		return dto.NewApiResponse(ErrDataBaseError, false)
	}

	if roleEntity.Name == role.Name && roleEntity.Comment == role.Description {
		service.logger.Infof("user %04d update role %s(%s) no change", role.Cid, role.Name, role.Description)
		return dto.NewApiResponse(dto.SuccessHandleRequest, true)
	}

	oldValue := fmt.Sprintf("%s(%s)", roleEntity.Name, roleEntity.Comment)

	roleEntity.Name = role.Name
	roleEntity.Comment = role.Description
	if err := service.repo.Save(roleEntity); err != nil {
		service.logger.Errorf("error occurred when update role: %v", err)
		return dto.NewApiResponse(ErrDataBaseError, false)
	}

	go func(role *DTO.UpdateRole, roleEntity *entity.Role, oldValue string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := service.auditLogClient.Log(ctx, &grpc.AuditLogRequest{
			Event:     entity.AuditEventRoleUpdated.Value,
			Subject:   fmt.Sprintf("%04d", role.Cid),
			Object:    strconv.Itoa(int(roleEntity.ID)),
			Ip:        role.Ip,
			UserAgent: role.UserAgent,
			OldValue:  oldValue,
			NewValue:  fmt.Sprintf("%s(%s)", roleEntity.Name, roleEntity.Comment),
		})
		if err != nil {
			service.logger.Errorf("error occurred when create audit log: %v", err)
		}
	}(role, roleEntity, oldValue)

	return dto.NewApiResponse(dto.SuccessHandleRequest, true)
}

var ErrRoleHasUsers = dto.NewApiStatus("ROLE_HAS_USERS", "角色下有用户", dto.HttpCodeConflict)

func (service *RoleService) Delete(role *DTO.DeleteRole) *dto.ApiResponse[bool] {
	perm := permission.Permission(role.Permission)
	if !perm.HasPermission(permission.RoleDelete) {
		service.logger.Errorf("user %04d no permission to delete role", role.Cid)
		return dto.NewApiResponse(dto.ErrNoPermission, false)
	}

	users, err := service.repo.GetRoleUsers(role.Id)
	if err != nil {
		service.logger.Errorf("error occurred when get role users: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.NewApiResponse(ErrRoleNotFound, false)
		}
		return dto.NewApiResponse(ErrDataBaseError, false)
	}

	if len(users) > 0 && !role.Force {
		service.logger.Errorf("role %d has users", role.Id)
		return dto.NewApiResponse(ErrRoleHasUsers, false)
	}

	if len(users) > 0 {
		err = service.repo.DeleteRole(role.Id)
	} else {
		err = service.repo.Delete(&entity.Role{ID: role.Id})
	}
	if err != nil {
		service.logger.Errorf("error occurred when delete role: %v", err)
		return dto.NewApiResponse(ErrDataBaseError, false)
	}

	go func(role *DTO.DeleteRole) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := service.auditLogClient.Log(ctx, &grpc.AuditLogRequest{
			Event:     entity.AuditEventRoleDeleted.Value,
			Subject:   fmt.Sprintf("%04d", role.Cid),
			Object:    strconv.Itoa(int(role.Id)),
			Ip:        role.Ip,
			UserAgent: role.UserAgent,
		})
		if err != nil {
			service.logger.Errorf("error occurred when create audit log: %v", err)
		}
	}(role)

	return dto.NewApiResponse(dto.SuccessHandleRequest, true)
}
