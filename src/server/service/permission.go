// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package service
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
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

type PermissionService struct {
	logger         logger.Interface
	userRepo       repository.UserInterface
	roleRepo       repository.RoleInterface
	emailClient    grpc.EmailClient
	auditLogClient grpc.AuditLogClient
}

func NewPermissionService(
	lg logger.Interface,
	userRepo repository.UserInterface,
	roleRepo repository.RoleInterface,
	emailClient grpc.EmailClient,
	auditLogClient grpc.AuditLogClient,
) *PermissionService {
	return &PermissionService{
		logger:         logger.NewLoggerAdapter(lg, "permission-service"),
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		emailClient:    emailClient,
		auditLogClient: auditLogClient,
	}
}

var (
	ErrPermissionNodeNotFound = dto.NewApiStatus("PERMISSION_NODE_NOT_FOUND", "权限节点不存在", dto.HttpCodeNotFound)
)

func checkDatabaseError[T comparable](err error) *dto.ApiResponse[T] {
	var zero T
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.NewApiResponse(ErrUserNotFound, zero)
	}
	return dto.NewApiResponse(dto.ErrServerError, zero)
}

func updatePermission(
	logger logger.Interface,
	perm permission.Permission,
	targetPerm uint64,
	changeData map[string]bool,
) (res *dto.ApiResponse[bool], resPerm permission.Permission, permGrant []string, permRevoke []string) {
	permGrant = make([]string, 0)
	permRevoke = make([]string, 0)
	resPerm = permission.Permission(targetPerm)

	for k, v := range changeData {
		if !permission.Permissions.IsValidEnum(k) {
			logger.Errorf("%s is not an valid permission node", k)
			res = dto.NewApiResponse(ErrPermissionNodeNotFound, false)
			return
		}
		node := permission.Permissions.GetEnum(k)
		if !perm.HasPermission(node.Data) {
			logger.Errorf("user has no permission on permission node %s", k)
			res = dto.NewApiResponse(dto.ErrNoPermission, false)
			return
		}
		if v {
			resPerm.Grant(node.Data)
			permGrant = append(permGrant, k)
		} else {
			resPerm.Revoke(node.Data)
			permRevoke = append(permRevoke, k)
		}
	}

	return
}

func (service *PermissionService) EditUserPermission(data *DTO.EditUserPermission) *dto.ApiResponse[bool] {
	perm := permission.Permission(data.Permission)
	if !perm.HasPermission(permission.UserEditPermission) {
		return dto.NewApiResponse(dto.ErrNoPermission, false)
	}

	targetUser, err := service.userRepo.GetById(data.UserId)
	if err != nil {
		service.logger.Errorf("get user failed: %v", err)
		return checkDatabaseError[bool](err)
	}

	res, targetPerm, permGrant, permRevoke := updatePermission(service.logger, perm, targetUser.Permission, data.Data)
	if res != nil {
		return res
	}

	if err := service.userRepo.Update(targetUser, map[string]interface{}{"permission": uint64(targetPerm)}); err != nil {
		service.logger.Errorf("update user permission failed: %v", err)
		return dto.NewApiResponse(dto.ErrServerError, false)
	}

	go func(data *DTO.EditUserPermission, targetUser *entity.User, grantList []string, revokeList []string) {
		user, _ := service.userRepo.GetById(data.Uid)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		auditLogRequest := &grpc.AuditLogRequest{
			Subject:   fmt.Sprintf("%04d", user.Cid),
			Object:    fmt.Sprintf("%04d", targetUser.Cid),
			Ip:        data.Ip,
			UserAgent: data.UserAgent,
		}
		if len(grantList) != 0 {
			auditLogRequest.Event = entity.AuditEventUserPermissionGrant.Value
			auditLogRequest.NewValue = strings.Join(grantList, ",")
			_, err := service.auditLogClient.Log(ctx, auditLogRequest)
			if err != nil {
				service.logger.Errorf("log user permission grant failed: %v", err)
			}
		}
		if len(revokeList) != 0 {
			auditLogRequest.Event = entity.AuditEventUserPermissionRevoke.Value
			auditLogRequest.NewValue = strings.Join(revokeList, ",")
			_, err := service.auditLogClient.Log(ctx, auditLogRequest)
			if err != nil {
				service.logger.Errorf("log user permission revoke failed: %v", err)
			}
		}
		_, err = service.emailClient.SendPermissionChange(ctx, &grpc.PermissionChange{
			TargetEmail: targetUser.Email,
			Cid:         auditLogRequest.Object,
			Permissions: "\n+" + strings.Join(grantList, "\n+") + "\n-" + strings.Join(revokeList, "\n-"),
			Operator:    auditLogRequest.Subject,
			Contact:     user.Email,
		})
	}(data, targetUser, permGrant, permRevoke)

	return dto.NewApiResponse(dto.SuccessHandleRequest, true)
}

func (service *PermissionService) EditRolePermission(data *DTO.EditRolePermission) *dto.ApiResponse[bool] {
	perm := permission.Permission(data.Permission)
	if !perm.HasPermission(permission.RoleEditPermission) {
		return dto.NewApiResponse(dto.ErrNoPermission, false)
	}

	targetRole, err := service.roleRepo.GetById(data.RoleId)
	if err != nil {
		service.logger.Errorf("get role failed: %v", err)
		return checkDatabaseError[bool](err)
	}

	res, targetPerm, permGrant, permRevoke := updatePermission(service.logger, perm, targetRole.Permission, data.Data)
	if res != nil {
		return res
	}

	if err := service.roleRepo.Update(targetRole, map[string]interface{}{"permission": uint64(targetPerm)}); err != nil {
		service.logger.Errorf("update role permission failed: %v", err)
		return dto.NewApiResponse(dto.ErrServerError, false)
	}

	go func(data *DTO.EditRolePermission, targetRole *entity.Role, grantList []string, revokeList []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		auditLogRequest := &grpc.AuditLogRequest{
			Subject:   fmt.Sprintf("%04d", data.Cid),
			Object:    fmt.Sprintf("%d: %s", targetRole.ID, targetRole.Name),
			Ip:        data.Ip,
			UserAgent: data.UserAgent,
		}
		if len(grantList) != 0 {
			auditLogRequest.Event = entity.AuditEventRolePermissionGrant.Value
			auditLogRequest.NewValue = strings.Join(grantList, ",")
			_, err := service.auditLogClient.Log(ctx, auditLogRequest)
			if err != nil {
				service.logger.Errorf("log role permission grant failed: %v", err)
			}
		}
		if len(revokeList) != 0 {
			auditLogRequest.Event = entity.AuditEventRolePermissionRevoke.Value
			auditLogRequest.NewValue = strings.Join(revokeList, ",")
			_, err := service.auditLogClient.Log(ctx, auditLogRequest)
			if err != nil {
				service.logger.Errorf("log role permission revoke failed: %v", err)
			}
		}
	}(data, targetRole, permGrant, permRevoke)

	return dto.NewApiResponse(dto.SuccessHandleRequest, true)
}

func (service *PermissionService) getUserAndRoles(userId uint, roleIds []uint) ([]*entity.Role, *entity.User, *dto.ApiResponse[bool]) {
	roles, err := service.roleRepo.GetByIds(roleIds)
	if err != nil {
		service.logger.Errorf("get role failed: %v", err)
		return nil, nil, checkDatabaseError[bool](err)
	}

	targetUser, err := service.userRepo.GetById(userId)
	if err != nil {
		service.logger.Errorf("get user failed: %v", err)
		return nil, nil, checkDatabaseError[bool](err)
	}
	return roles, targetUser, nil
}

func (service *PermissionService) GrantUserRole(data *DTO.GrantUserRole) *dto.ApiResponse[bool] {
	perm := permission.Permission(data.Permission)
	if !perm.HasPermission(permission.UserEditRole) {
		return dto.NewApiResponse(dto.ErrNoPermission, false)
	}

	roles, targetUser, res := service.getUserAndRoles(data.UserId, data.RoleIds)
	if res != nil {
		return res
	}

	if err := service.userRepo.GrantRole(targetUser.ID, data.RoleIds); err != nil {
		service.logger.Errorf("grant user role failed: %v", err)
		return checkDatabaseError[bool](err)
	}

	go func(data *DTO.GrantUserRole, targetUser *entity.User, roles []*entity.Role) {
		newRoles := make([]string, len(roles))
		utils.ForEach(roles, func(index int, role *entity.Role) {
			newRoles[index] = role.Name
		})
		user, _ := service.userRepo.GetById(data.Uid)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := service.auditLogClient.Log(ctx, &grpc.AuditLogRequest{
			Event:     entity.AuditEventRoleGrant.Value,
			Subject:   fmt.Sprintf("%04d", user.Cid),
			Object:    fmt.Sprintf("%04d", targetUser.Cid),
			Ip:        data.Ip,
			UserAgent: data.UserAgent,
			NewValue:  strings.Join(newRoles, ","),
		})
		if err != nil {
			service.logger.Errorf("log role grant failed: %v", err)
		}
		_, err = service.emailClient.SendRoleChange(ctx, &grpc.RoleChange{
			TargetEmail: []string{targetUser.Email},
			Cid:         fmt.Sprintf("%04d", targetUser.Cid),
			Roles:       "\n+" + strings.Join(newRoles, "\n+"),
			Operator:    fmt.Sprintf("%04d", user.Cid),
			Contact:     user.Email,
		})
		if err != nil {
			service.logger.Errorf("send role change email failed: %v", err)
		}
	}(data, targetUser, roles)

	return dto.NewApiResponse(dto.SuccessHandleRequest, true)
}

func (service *PermissionService) RevokeUserRole(data *DTO.RevokeUserRole) *dto.ApiResponse[bool] {
	perm := permission.Permission(data.Permission)
	if !perm.HasPermission(permission.UserEditRole) {
		return dto.NewApiResponse(dto.ErrNoPermission, false)
	}

	roles, targetUser, res := service.getUserAndRoles(data.UserId, data.RoleIds)
	if res != nil {
		return res
	}

	if err := service.userRepo.RevokeRole(targetUser.ID, data.RoleIds); err != nil {
		service.logger.Errorf("revoke user role failed: %v", err)
		return checkDatabaseError[bool](err)
	}

	go func(data *DTO.RevokeUserRole, targetUser *entity.User, roles []*entity.Role) {
		oldRoles := make([]string, len(roles))
		utils.ForEach(roles, func(index int, role *entity.Role) {
			oldRoles[index] = role.Name
		})
		user, _ := service.userRepo.GetById(data.Uid)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := service.auditLogClient.Log(ctx, &grpc.AuditLogRequest{
			Event:     entity.AuditEventRoleRevoke.Value,
			Subject:   fmt.Sprintf("%04d", user.Cid),
			Object:    fmt.Sprintf("%04d", targetUser.Cid),
			Ip:        data.Ip,
			UserAgent: data.UserAgent,
			OldValue:  strings.Join(oldRoles, ","),
		})
		if err != nil {
			service.logger.Errorf("log role revoke failed: %v", err)
		}
		_, err = service.emailClient.SendRoleChange(ctx, &grpc.RoleChange{
			TargetEmail: []string{targetUser.Email},
			Cid:         fmt.Sprintf("%04d", targetUser.Cid),
			Roles:       "\n-" + strings.Join(oldRoles, "\n-"),
			Operator:    fmt.Sprintf("%04d", user.Cid),
			Contact:     user.Email,
		})
		if err != nil {
			service.logger.Errorf("send role change email failed: %v", err)
		}
	}(data, targetUser, roles)

	return dto.NewApiResponse(dto.SuccessHandleRequest, true)
}

func (service *PermissionService) getRoleAndUsers(roleId uint, userIds []uint) ([]*entity.User, *entity.Role, *dto.ApiResponse[bool]) {
	users, err := service.userRepo.GetByIds(userIds)
	if err != nil {
		service.logger.Errorf("get user failed: %v", err)
		return nil, nil, checkDatabaseError[bool](err)
	}

	targetRole, err := service.roleRepo.GetById(roleId)
	if err != nil {
		service.logger.Errorf("get role failed: %v", err)
		return nil, nil, checkDatabaseError[bool](err)
	}

	return users, targetRole, nil
}
func (service *PermissionService) GrantRoleUser(data *DTO.GrantRoleUser) *dto.ApiResponse[bool] {
	perm := permission.Permission(data.Permission)
	if !perm.HasPermission(permission.UserEditRole) {
		return dto.NewApiResponse(dto.ErrNoPermission, false)
	}

	users, role, res := service.getRoleAndUsers(data.RoleId, data.UserIds)
	if res != nil {
		return res
	}

	if err := service.roleRepo.GrantUser(role.ID, data.UserIds); err != nil {
		service.logger.Errorf("grant role user failed: %v", err)
		return checkDatabaseError[bool](err)
	}

	go func(data *DTO.GrantRoleUser, users []*entity.User, role *entity.Role) {
		userCids := make([]string, len(users))
		utils.ForEach(users, func(index int, user *entity.User) {
			userCids[index] = fmt.Sprintf("%04d", user.Cid)
		})

		user, _ := service.userRepo.GetById(data.Uid)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(len(users)+1)*5*time.Second)
		defer cancel()

		_, err := service.auditLogClient.Log(ctx, &grpc.AuditLogRequest{
			Event:     entity.AuditEventRoleGrant.Value,
			Subject:   fmt.Sprintf("%04d", data.Cid),
			Object:    fmt.Sprintf("%s: %d", role.Name, role.ID),
			Ip:        data.Ip,
			UserAgent: data.UserAgent,
			NewValue:  strings.Join(userCids, ","),
		})
		if err != nil {
			service.logger.Errorf("log role grant failed: %v", err)
		}

		emails := make([]string, len(users))
		utils.ForEach(users, func(index int, user *entity.User) {
			emails[index] = user.Email
		})
		_, err = service.emailClient.SendRoleChange(ctx, &grpc.RoleChange{
			TargetEmail: emails,
			Cid:         strings.Join(userCids, ","),
			Roles:       fmt.Sprintf("+%s", role.Name),
			Operator:    fmt.Sprintf("%04d", user.Cid),
			Contact:     user.Email,
		})
		if err != nil {
			service.logger.Errorf("send role change email failed: %v", err)
		}
	}(data, users, role)

	return dto.NewApiResponse(dto.SuccessHandleRequest, true)
}

func (service *PermissionService) RevokeRoleUser(data *DTO.RevokeRoleUser) *dto.ApiResponse[bool] {
	perm := permission.Permission(data.Permission)
	if !perm.HasPermission(permission.UserEditRole) {
		return dto.NewApiResponse(dto.ErrNoPermission, false)
	}

	users, role, res := service.getRoleAndUsers(data.RoleId, data.UserIds)
	if res != nil {
		return res
	}

	if err := service.roleRepo.RevokeUser(role.ID, data.UserIds); err != nil {
		service.logger.Errorf("revoke role user failed: %v", err)
		return checkDatabaseError[bool](err)
	}

	go func(data *DTO.RevokeRoleUser, users []*entity.User, role *entity.Role) {
		userCids := make([]string, len(users))
		utils.ForEach(users, func(index int, user *entity.User) {
			userCids[index] = fmt.Sprintf("%04d", user.Cid)
		})

		user, _ := service.userRepo.GetById(data.Uid)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(len(users)+1)*5*time.Second)
		defer cancel()

		_, err := service.auditLogClient.Log(ctx, &grpc.AuditLogRequest{
			Event:     entity.AuditEventRoleRevoke.Value,
			Subject:   fmt.Sprintf("%04d", data.Cid),
			Object:    fmt.Sprintf("%s: %d", role.Name, role.ID),
			Ip:        data.Ip,
			UserAgent: data.UserAgent,
			OldValue:  strings.Join(userCids, ","),
		})
		if err != nil {
			service.logger.Errorf("log role revoke failed: %v", err)
		}

		emails := make([]string, len(users))
		utils.ForEach(users, func(index int, user *entity.User) {
			emails[index] = user.Email
		})
		_, err = service.emailClient.SendRoleChange(ctx, &grpc.RoleChange{
			TargetEmail: emails,
			Cid:         strings.Join(userCids, ","),
			Roles:       fmt.Sprintf("-%s", role.Name),
			Operator:    fmt.Sprintf("%04d", user.Cid),
			Contact:     user.Email,
		})
		if err != nil {
			service.logger.Errorf("send role change email failed: %v", err)
		}
	}(data, users, role)

	return dto.NewApiResponse(dto.SuccessHandleRequest, true)
}
