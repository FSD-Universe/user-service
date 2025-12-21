// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package service
package service

import (
	"context"
	"fmt"
	"user-service/src/interfaces/global"
	"user-service/src/interfaces/grpc"
	pb "user-service/src/interfaces/grpc"
	"user-service/src/interfaces/repository"
	DTO "user-service/src/interfaces/server/dto"

	"half-nothing.cn/service-core/interfaces/database/entity"
	"half-nothing.cn/service-core/interfaces/http/dto"
	"half-nothing.cn/service-core/interfaces/logger"
	"half-nothing.cn/service-core/utils"
)

type UserService struct {
	logger         logger.Interface
	repo           repository.UserInterface
	emailClient    grpc.EmailClient
	auditLogClient grpc.AuditLogClient
}

func NewUserService(
	lg logger.Interface,
	repo repository.UserInterface,
	emailClient grpc.EmailClient,
	auditLogClient grpc.AuditLogClient,
) *UserService {
	return &UserService{
		logger:         logger.NewLoggerAdapter(lg, "user-service"),
		repo:           repo,
		emailClient:    emailClient,
		auditLogClient: auditLogClient,
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
)

func (u *UserService) Register(form *DTO.UserRegister) *dto.ApiResponse[bool] {
	code, err := u.emailClient.VerifyEmailCode(context.Background(), &pb.VerifyCode{
		Email: form.Email,
		Code:  form.Code,
		Cid:   int32(form.Cid),
	})
	if err != nil {
		u.logger.Errorf("error occurred when verify email code: %v", err)
		return nil
	}

	switch code.Code {
	case CodeExpired:
		u.logger.Errorf("email code %s for %s(%04d) expired", form.Code, form.Email, form.Cid)
		return dto.NewApiResponse(ErrCodeExpired, false)
	case CodeInvalid:
		u.logger.Errorf("email code %s for %s(%04d) invalid", form.Code, form.Email, form.Cid)
		return dto.NewApiResponse(ErrCodeInvalid, false)
	case CodeError:
		u.logger.Errorf("email code %s for %s(%04d) error", form.Code, form.Email, form.Cid)
		return dto.NewApiResponse(ErrCodeUnknown, false)
	case CodeValid:
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

	go func() {
		_, err := u.emailClient.SendWelcome(context.Background(), &pb.Welcome{
			TargetEmail: form.Email,
			Cid:         fmt.Sprintf("%04d", form.Cid),
		})
		if err != nil {
			u.logger.Errorf("error occurred when send welcome email: %v", err)
		}
		_, err = u.auditLogClient.Log(context.Background(), &pb.AuditLogRequest{
			Event:     entity.AuditEventUserRegistered.Value,
			Subject:   fmt.Sprintf("%04d", form.Cid),
			Object:    fmt.Sprintf("%s(%s)", user.Username, user.Email),
			Ip:        form.Ip,
			UserAgent: form.UserAgent,
		})
		if err != nil {
			u.logger.Errorf("error occurred when log audit: %v", err)
		}
	}()

	return dto.NewApiResponse(dto.SuccessHandleRequest, true)
}
