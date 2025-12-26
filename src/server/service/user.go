// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package service
package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"user-service/src/interfaces/content"
	"user-service/src/interfaces/global"
	"user-service/src/interfaces/grpc"
	pb "user-service/src/interfaces/grpc"
	"user-service/src/interfaces/repository"
	DTO "user-service/src/interfaces/server/dto"

	"gorm.io/gorm"
	"half-nothing.cn/service-core/interfaces/database/entity"
	"half-nothing.cn/service-core/interfaces/http/dto"
	"half-nothing.cn/service-core/interfaces/logger"
	"half-nothing.cn/service-core/permission"
	"half-nothing.cn/service-core/utils"
)

type UserService struct {
	logger logger.Interface
	repo   repository.UserInterface
	client *content.GrpcClientManager
}

func NewUserService(
	lg logger.Interface,
	repo repository.UserInterface,
	client *content.GrpcClientManager,
) *UserService {
	return &UserService{
		logger: logger.NewLoggerAdapter(lg, "user-service"),
		repo:   repo,
		client: client,
	}
}

const (
	CodeValid int32 = iota
	CodeExpired
	CodeInvalid
	CodeError
)

var (
	ErrCodeExpired     = dto.NewApiStatus("CODE_EXPIRED", "验证码已过期", dto.HttpCodeBadRequest)
	ErrCodeInvalid     = dto.NewApiStatus("CODE_INVALID", "验证码错误", dto.HttpCodeBadRequest)
	ErrCodeUnknown     = dto.NewApiStatus("CODE_UNKNOWN", "服务器错误", dto.HttpCodeInternalError)
	ErrPasswordEncrypt = dto.NewApiStatus("PASSWORD_ENCRYPT_ERROR", "密码加密错误", dto.HttpCodeInternalError)
	ErrDataBaseError   = dto.NewApiStatus("DATABASE_ERROR", "数据库错误", dto.HttpCodeInternalError)
	ErrRegistered      = dto.NewApiStatus("USER_REGISTERED", "邮箱、用户名或呼号已被注册", dto.HttpCodeBadRequest)
)

func verifyEmailCode[T comparable](u *UserService, email string, code string) *dto.ApiResponse[T] {
	var zero T
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := u.client.EmailClient().VerifyEmailCode(ctx, &pb.VerifyCode{Email: email, Code: code})
	if err != nil {
		u.logger.Errorf("error occurred when verify email code: %v", err)
		return dto.NewApiResponse[T](dto.ErrServerError, zero)
	}
	switch res.Code {
	case CodeExpired:
		return dto.NewApiResponse[T](ErrCodeExpired, zero)
	case CodeInvalid:
		return dto.NewApiResponse[T](ErrCodeInvalid, zero)
	case CodeError:
		return dto.NewApiResponse[T](ErrCodeUnknown, zero)
	case CodeValid:
		fallthrough
	default:
		return nil
	}
}

func (u *UserService) removeEmailCode(ctx context.Context, email string) {
	_, err := u.client.EmailClient().RemoveEmailCode(ctx, &pb.RemoveVerifyCode{Email: email})
	if err != nil {
		u.logger.Errorf("error occurred when remove email code: %v", err)
	}
}

func (u *UserService) Register(form *DTO.UserRegister) *dto.ApiResponse[bool] {
	if exist, err := u.repo.CheckCidUsernameAndEmail(uint(form.Cid), form.Username, form.Email); !exist {
		u.logger.Errorf("error occurred when check cid username and email: %v", err)
		return dto.NewApiResponse(ErrRegistered, false)
	}

	if res := verifyEmailCode[bool](u, form.Email, form.Code); res != nil {
		return res
	}

	hashedPassword, err := utils.BcryptEncrypt([]byte(form.Password), *global.BcryptCost)
	if err != nil {
		u.logger.Errorf("error occurred when encrypt password: %v", err)
		return dto.NewApiResponse(ErrPasswordEncrypt, false)
	}

	user := &entity.User{
		Username: form.Username,
		Email:    form.Email,
		Cid:      uint(form.Cid),
		Password: string(hashedPassword),
	}
	if err := u.repo.Save(user); err != nil {
		u.logger.Errorf("error occurred when save user: %v", err)
		return dto.NewApiResponse(ErrDataBaseError, false)
	}

	go func(u *UserService, form *DTO.UserRegister, user *entity.User) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_, err := u.client.EmailClient().SendWelcome(ctx, &pb.Welcome{
			TargetEmail: form.Email,
			Cid:         fmt.Sprintf("%04d", form.Cid),
		})
		if err != nil {
			u.logger.Errorf("error occurred when send welcome email: %v", err)
		}
		_, err = u.client.AuditLogClient().Log(ctx, &pb.AuditLogRequest{
			Event:     entity.AuditEventUserRegistered.Value,
			Subject:   fmt.Sprintf("%04d", form.Cid),
			Object:    fmt.Sprintf("%s(%s)", user.Username, user.Email),
			Ip:        form.Ip,
			UserAgent: form.UserAgent,
		})
		if err != nil {
			u.logger.Errorf("error occurred when log audit: %v", err)
		}
		u.removeEmailCode(ctx, form.Email)
	}(u, form, user)

	return dto.NewApiResponse(dto.SuccessHandleRequest, true)
}

func (u *UserService) CheckAvailability(form *DTO.UserCheckAvailability) *dto.ApiResponse[bool] {
	exist, err := u.repo.CheckCidUsernameAndEmail(uint(form.Cid), form.Username, form.Email)
	if err != nil {
		u.logger.Errorf("error occurred when check cid username and email: %v", err)
	}
	return dto.NewApiResponse(dto.SuccessHandleRequest, exist)
}

var (
	ErrUserNotExist = dto.NewApiStatus("USER_NOT_EXIST", "用户不存在", dto.HttpCodeNotFound)
)

func (u *UserService) ResetPassword(form *DTO.UserResetPassword) *dto.ApiResponse[bool] {
	if res := verifyEmailCode[bool](u, form.Email, form.Code); res != nil {
		return res
	}

	user, err := u.repo.GetByUsernameOrEmail(form.Email)
	if err != nil {
		u.logger.Errorf("ResetPassword handle fail, get user err, %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.NewApiResponse(ErrUserNotExist, false)
		}
		return dto.NewApiResponse(ErrDataBaseError, false)
	}

	hashedPassword, err := utils.BcryptEncrypt([]byte(form.Password), *global.BcryptCost)
	if err != nil {
		u.logger.Errorf("error occurred when encrypt password: %v", err)
		return dto.NewApiResponse(ErrPasswordEncrypt, false)
	}

	updates := map[string]interface{}{
		"password": string(hashedPassword),
	}

	if err := u.repo.Update(user, updates); err != nil {
		u.logger.Errorf("ResetPassword handle fail, save user err, %v", err)
		return dto.NewApiResponse(ErrDataBaseError, false)
	}

	go func(u *UserService, form *DTO.UserResetPassword, user *entity.User) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := u.client.AuditLogClient().Log(ctx, &grpc.AuditLogRequest{
			Event:     entity.AuditEventUserResetPassword.Value,
			Subject:   fmt.Sprintf("%04d", user.Cid),
			Object:    fmt.Sprintf("%s(%s)", user.Username, user.Email),
			Ip:        form.Ip,
			UserAgent: form.UserAgent,
		})
		if err != nil {
			u.logger.Errorf("error occurred when log audit: %v", err)
		}
		_, err = u.client.EmailClient().SendPasswordReset(ctx, &pb.PasswordReset{
			TargetEmail: form.Email,
			Cid:         fmt.Sprintf("%04d", user.Cid),
			Time:        time.Now().Format(time.RFC3339),
			Ip:          form.Ip,
			UserAgent:   form.UserAgent,
		})
		if err != nil {
			u.logger.Errorf("error occurred when send password reset email: %v", err)
		}
		u.removeEmailCode(ctx, form.Email)
	}(u, form, user)

	return dto.NewApiResponse(dto.SuccessHandleRequest, true)
}

func (u *UserService) GetPages(page *DTO.GetUserPage) *dto.ApiResponse[*DTO.GetUserPageResponse] {
	perm := permission.Permission(page.Permission)
	if !perm.HasPermission(permission.UserShowList) {
		u.logger.Errorf("user %04d no permission to show user list", page.Cid)
		return dto.NewApiResponse[*DTO.GetUserPageResponse](dto.ErrNoPermission, nil)
	}
	users, total, err := u.repo.GetPages(page.PageNum, page.PageSize, page.Search)
	if err != nil {
		u.logger.Errorf("error occurred when get pages: %v", err)
		return dto.NewApiResponse[*DTO.GetUserPageResponse](ErrDataBaseError, nil)
	}
	userInfos := make([]*DTO.FullUserInfo, len(users))
	utils.ForEach(users, func(index int, element *entity.User) {
		userInfos[index] = &DTO.FullUserInfo{}
		userInfos[index].FromUserEntity(element)
	})
	return dto.NewApiResponse(dto.SuccessHandleRequest, &DTO.GetUserPageResponse{
		Data:     userInfos,
		Total:    int(total),
		PageNum:  page.PageNum,
		PageSize: page.PageSize,
	})
}

var (
	ErrUserNotFound = dto.NewApiStatus("USER_NOT_FOUND", "用户不存在", dto.HttpCodeNotFound)
)

func (u *UserService) GetSelfData(data *DTO.GetCurrentUserData) *dto.ApiResponse[*DTO.UserInfo] {
	user, err := u.repo.GetById(data.Uid)
	if err != nil {
		u.logger.Errorf("GetSelfData handle fail, get user err, %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.NewApiResponse[*DTO.UserInfo](ErrUserNotFound, nil)
		}
		return dto.NewApiResponse[*DTO.UserInfo](ErrDataBaseError, nil)
	}
	userInfo := &DTO.UserInfo{}
	userInfo.FromUserEntity(user)
	return dto.NewApiResponse(dto.SuccessHandleRequest, userInfo)
}

func (u *UserService) GetData(data *DTO.GetUserData) *dto.ApiResponse[*DTO.FullUserInfo] {
	perm := permission.Permission(data.Permission)
	if !perm.HasPermission(permission.UserShowList) {
		u.logger.Errorf("user %04d no permission to get user data", data.Cid)
		return dto.NewApiResponse[*DTO.FullUserInfo](dto.ErrNoPermission, nil)
	}
	user, err := u.repo.GetById(data.Id)
	if err != nil {
		u.logger.Errorf("GetData handle fail, get user err, %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.NewApiResponse[*DTO.FullUserInfo](ErrUserNotFound, nil)
		}
		return dto.NewApiResponse[*DTO.FullUserInfo](ErrDataBaseError, nil)
	}
	userInfo := &DTO.FullUserInfo{}
	userInfo.FromUserEntity(user)
	return dto.NewApiResponse(dto.SuccessHandleRequest, userInfo)
}

func (u *UserService) UpdateSelfData(data *DTO.UpdateCurrentUserData) *dto.ApiResponse[*DTO.UserInfo] {
	user, err := u.repo.GetById(data.Uid)
	if err != nil {
		u.logger.Errorf("UpdateSelfData handle fail, get user err, %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.NewApiResponse[*DTO.UserInfo](ErrUserNotFound, nil)
		}
		return dto.NewApiResponse[*DTO.UserInfo](ErrDataBaseError, nil)
	}

	var oldEmail string
	updates := map[string]interface{}{}
	if data.Username != "" && user.Username != data.Username {
		updates["username"] = data.Username
	}
	if data.Email != "" && user.Email != data.Email {
		if res := verifyEmailCode[*DTO.UserInfo](u, data.Email, data.EmailCode); res != nil {
			return res
		}
		oldEmail = user.Email
		updates["email"] = data.Email
	}
	if data.QQ != "" && (user.QQ == nil || *user.QQ != data.QQ) {
		updates["qq"] = &data.QQ
	}
	if (data.ImageId != nil && user.ImageId != nil && *user.ImageId != *data.ImageId) ||
		(data.ImageId == nil && user.ImageId != nil) ||
		(data.ImageId != nil && user.ImageId == nil) {
		updates["image_id"] = data.ImageId
	}

	if len(updates) == 0 {
		return dto.NewApiResponse[*DTO.UserInfo](dto.ErrErrorParam, nil)
	}

	if err := u.repo.Update(user, updates); err != nil {
		u.logger.Errorf("UpdateSelfData handle fail, save user err, %v", err)
		return dto.NewApiResponse[*DTO.UserInfo](ErrDataBaseError, nil)
	}

	user, err = u.repo.GetById(data.Uid)
	if err != nil {
		u.logger.Errorf("UpdateSelfData handle fail, get user err, %v", err)
		return dto.NewApiResponse[*DTO.UserInfo](ErrDataBaseError, nil)
	}

	if oldEmail != "" {
		go func(u *UserService, user *entity.User, oldEmail string) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := u.client.EmailClient().SendEmailChange(ctx, &grpc.EmailChange{
				TargetEmail: oldEmail,
				Cid:         fmt.Sprintf("%04d", user.Cid),
				Email:       user.Email,
				Time:        time.Now().Format(time.RFC3339),
				Ip:          data.Ip,
				UserAgent:   data.UserAgent,
			})
			if err != nil {
				u.logger.Errorf("error occurred when send email change email: %v", err)
			}
		}(u, user, oldEmail)
	}

	userInfo := &DTO.UserInfo{}
	userInfo.FromUserEntity(user)

	return dto.NewApiResponse(dto.SuccessHandleRequest, userInfo)
}

func (u *UserService) UpdateData(data *DTO.UpdateUserData) *dto.ApiResponse[bool] {
	perm := permission.Permission(data.Permission)
	if !perm.HasPermission(permission.UserEditInfo) {
		u.logger.Errorf("user %04d no permission to update user data", data.Cid)
		return dto.NewApiResponse[bool](dto.ErrNoPermission, false)
	}
	user, err := u.repo.GetById(data.Id)
	if err != nil {
		u.logger.Errorf("UpdateData handle fail, get user err, %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.NewApiResponse[bool](ErrUserNotFound, false)
		}
		return dto.NewApiResponse[bool](ErrDataBaseError, false)
	}

	oldValue := map[string]interface{}{}
	updates := map[string]interface{}{}
	if data.Username != "" && user.Username != data.Username {
		oldValue["username"] = user.Username
		updates["username"] = data.Username
	}
	if data.Email != "" && user.Email != data.Email {
		oldValue["email"] = user.Email
		updates["email"] = data.Email
	}
	if data.QQ != "" && (user.QQ == nil || *user.QQ != data.QQ) {
		oldValue["qq"] = user.QQ
		updates["qq"] = &data.QQ
	}
	if data.Password != "" {
		password, err := utils.BcryptEncrypt([]byte(data.Password), *global.BcryptCost)
		if err != nil {
			u.logger.Errorf("UpdateData handle fail, bcrypt error, %v", err)
			return dto.NewApiResponse[bool](ErrPasswordEncrypt, false)
		}
		updates["password"] = string(password)
	}

	if err := u.repo.Update(user, updates); err != nil {
		u.logger.Errorf("UpdateData handle fail, save user err, %v", err)
		return dto.NewApiResponse[bool](ErrDataBaseError, false)
	}

	go func(u *UserService, data *DTO.UpdateUserData, user *entity.User, oldValue map[string]interface{}, newValue map[string]interface{}) {
		oldValueStr, _ := json.Marshal(oldValue)
		newValueStr, _ := json.Marshal(newValue)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_, err = u.client.AuditLogClient().Log(ctx, &pb.AuditLogRequest{
			Event:     entity.AuditEventUserInformationEdit.Value,
			Subject:   fmt.Sprintf("%04d", data.Cid),
			Object:    fmt.Sprintf("%04d", user.Cid),
			Ip:        data.Ip,
			UserAgent: data.UserAgent,
			OldValue:  string(oldValueStr),
			NewValue:  string(newValueStr),
		})
		if err != nil {
			u.logger.Errorf("error occurred when log audit: %v", err)
		}
	}(u, data, user, oldValue, updates)

	return dto.NewApiResponse(dto.SuccessHandleRequest, true)
}

var (
	ErrOldPassword = dto.NewApiStatus("OLD_PASSWORD_ERROR", "原密码错误", dto.HttpCodeBadRequest)
)

func (u *UserService) UpdatePassword(data *DTO.UpdateUserPassword) *dto.ApiResponse[bool] {
	user, err := u.repo.GetById(data.Uid)
	if err != nil {
		u.logger.Errorf("UpdatePassword handle fail, get user err, %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.NewApiResponse[bool](ErrUserNotFound, false)
		}
		return dto.NewApiResponse[bool](ErrDataBaseError, false)
	}

	if !utils.BcryptCompare([]byte(data.OldPassword), []byte(user.Password)) {
		u.logger.Errorf("UpdatePassword handle fail, old password error")
		return dto.NewApiResponse[bool](ErrOldPassword, false)
	}

	password, err := utils.BcryptEncrypt([]byte(data.NewPassword), *global.BcryptCost)
	if err != nil {
		u.logger.Errorf("UpdatePassword handle fail, bcrypt error, %v", err)
		return dto.NewApiResponse[bool](ErrPasswordEncrypt, false)
	}

	updates := map[string]interface{}{
		"password": string(password),
	}

	if err := u.repo.Update(user, updates); err != nil {
		u.logger.Errorf("UpdatePassword handle fail, save user err, %v", err)
		return dto.NewApiResponse[bool](ErrDataBaseError, false)
	}

	go func(u *UserService, data *DTO.UpdateUserPassword, user *entity.User) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := u.client.EmailClient().SendPasswordChange(ctx, &pb.PasswordChange{
			TargetEmail: user.Email,
			Cid:         fmt.Sprintf("%04d", user.Cid),
			Time:        time.Now().Format(time.RFC3339),
			Ip:          data.Ip,
			UserAgent:   data.UserAgent,
		})
		if err != nil {
			u.logger.Errorf("error occurred when send password change email: %v", err)
		}
	}(u, data, user)

	return dto.NewApiResponse(dto.SuccessHandleRequest, true)
}

func (u *UserService) Ban(data *DTO.BanUser) *dto.ApiResponse[bool] {
	perm := permission.Permission(data.Permission)
	if !perm.HasPermission(permission.UserBan) {
		u.logger.Errorf("user %04d no permission to ban user", data.Cid)
		return dto.NewApiResponse[bool](dto.ErrNoPermission, false)
	}

	user, err := u.repo.GetById(data.Id)
	if err != nil {
		u.logger.Errorf("Ban handle fail, get user err, %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.NewApiResponse[bool](ErrUserNotFound, false)
		}
		return dto.NewApiResponse[bool](ErrDataBaseError, false)
	}

	var bannedUntil sql.NullTime
	bannedUntil.Valid = data.BannedSeconds > 0
	bannedUntil.Time = time.Now().Add(time.Duration(data.BannedSeconds) * time.Second)

	if err := u.repo.Ban(user.ID, bannedUntil); err != nil {
		return dto.NewApiResponse[bool](ErrDataBaseError, false)
	}

	go func(data *DTO.BanUser, user *entity.User, bannedUntil sql.NullTime) {
		operator, _ := u.repo.GetById(data.Uid)
		var bannedTime string
		if !bannedUntil.Valid {
			bannedTime = "永不解封"
		} else {
			bannedTime = bannedUntil.Time.Format(time.RFC3339)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := u.client.AuditLogClient().Log(ctx, &pb.AuditLogRequest{
			Event:     entity.AuditEventUserBan.Value,
			Subject:   fmt.Sprintf("%04d", operator.Cid),
			Object:    fmt.Sprintf("%04d", user.Cid),
			Ip:        data.Ip,
			UserAgent: data.UserAgent,
			NewValue:  data.Reason,
		})
		if err != nil {
			u.logger.Errorf("error occurred when log audit: %v", err)
		}
		_, err = u.client.EmailClient().SendBanned(ctx, &pb.Banned{
			TargetEmail: user.Email,
			Cid:         fmt.Sprintf("%04d", user.Cid),
			Reason:      data.Reason,
			Time:        bannedTime,
			Operator:    fmt.Sprintf("%04d", operator.Cid),
			Contact:     operator.Email,
		})
		if err != nil {
			u.logger.Errorf("error occurred when send banned email: %v", err)
		}
	}(data, user, bannedUntil)

	return dto.NewApiResponse(dto.SuccessHandleRequest, true)
}

func (u *UserService) Unban(data *DTO.UnbanUser) *dto.ApiResponse[bool] {
	perm := permission.Permission(data.Permission)
	if !perm.HasPermission(permission.UserBan) {
		u.logger.Errorf("user %04d no permission to unban user", data.Cid)
		return dto.NewApiResponse[bool](dto.ErrNoPermission, false)
	}

	user, err := u.repo.GetById(data.Id)
	if err != nil {
		u.logger.Errorf("Unban handle fail, get user err, %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.NewApiResponse[bool](ErrUserNotFound, false)
		}
		return dto.NewApiResponse[bool](ErrDataBaseError, false)
	}

	if err := u.repo.Unban(user.ID); err != nil {
		return dto.NewApiResponse[bool](ErrDataBaseError, false)
	}

	go func(data *DTO.UnbanUser, user *entity.User) {
		operator, _ := u.repo.GetById(data.Uid)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := u.client.AuditLogClient().Log(ctx, &pb.AuditLogRequest{
			Event:     entity.AuditEventUserUnban.Value,
			Subject:   fmt.Sprintf("%04d", operator.Cid),
			Object:    fmt.Sprintf("%04d", user.Cid),
			Ip:        data.Ip,
			UserAgent: data.UserAgent,
		})
		if err != nil {
			u.logger.Errorf("error occurred when log audit: %v", err)
		}
		_, err = u.client.EmailClient().SendUnbanned(ctx, &pb.Unbanned{
			TargetEmail: user.Email,
			Cid:         fmt.Sprintf("%04d", user.Cid),
			Operator:    fmt.Sprintf("%04d", operator.Cid),
			Contact:     operator.Email,
		})
		if err != nil {
			u.logger.Errorf("error occurred when send unbanned email: %v", err)
		}
	}(data, user)

	return dto.NewApiResponse(dto.SuccessHandleRequest, true)
}
