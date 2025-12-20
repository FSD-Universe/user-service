// Copyright (c) 2025 Half_nothing
// SPDX-License-Identifier: MIT

// Package service
package service

import (
	"errors"
	"fmt"
	"time"
	"user-service/src/interfaces/repository"
	DTO "user-service/src/interfaces/server/dto"
	"user-service/src/interfaces/server/service"

	"golang.org/x/crypto/bcrypt"
	"half-nothing.cn/service-core/interfaces/http/dto"
	"half-nothing.cn/service-core/interfaces/http/jwt"
	"half-nothing.cn/service-core/interfaces/logger"
)

type AuthService struct {
	logger       logger.Interface
	userRepo     repository.UserInterface
	claimFactory jwt.ClaimFactoryInterface
}

func NewAuthService(
	lg logger.Interface,
	userRepo repository.UserInterface,
	claimFactory jwt.ClaimFactoryInterface,
) *AuthService {
	return &AuthService{
		logger:       logger.NewLoggerAdapter(lg, "user-service"),
		userRepo:     userRepo,
		claimFactory: claimFactory,
	}
}

func (s *AuthService) Login(form *DTO.UserLogin) *dto.ApiResponse[*DTO.UserLoginResponse] {
	userId := repository.GetUserId(form.Username)
	user, err := userId.GetUser(s.userRepo)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			s.logger.Errorf("UserLogin handle fail, %s user not found", form.Username)
			return dto.NewApiResponse[*DTO.UserLoginResponse](service.ErrUsernameOrPasswordError, nil)
		}
		s.logger.Errorf("UserLogin handle fail, get user err, %v", err)
		return dto.NewApiResponse[*DTO.UserLoginResponse](dto.ErrServerError, nil)
	}

	updates := make(map[string]interface{})

	if user.Banned && user.BannedUntil.Valid && user.BannedUntil.Time.Before(time.Now()) {
		user.Banned = false
		user.BannedUntil.Valid = false
		updates["banned"] = user.Banned
		updates["banned_until"] = user.BannedUntil
	}

	if user.Banned {
		if user.BannedUntil.Valid {
			return dto.NewApiResponse[*DTO.UserLoginResponse](
				dto.NewApiStatus(
					"USER_BANNED",
					fmt.Sprintf("您已被封禁，解封时间：%s", user.BannedUntil.Time.Format("2006-01-02 15:04:05")),
					dto.HttpCodePermissionDenied,
				),
				nil,
			)
		}

		return dto.NewApiResponse[*DTO.UserLoginResponse](service.ErrUserBanned, nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		return dto.NewApiResponse[*DTO.UserLoginResponse](service.ErrUsernameOrPasswordError, nil)
	}

	user.LastLoginTime.Valid = true
	user.LastLoginTime.Time = time.Now()
	user.LastLoginIP = &form.Ip
	updates["last_login_time"] = user.LastLoginTime
	updates["last_login_ip"] = user.LastLoginIP

	if err := s.userRepo.Update(user, updates); err != nil {
		s.logger.Errorf("UserLogin handle fail, save user err, %v", err)
		return dto.NewApiResponse[*DTO.UserLoginResponse](dto.ErrServerError, nil)
	}

	token, err := s.claimFactory.GenerateKey(s.claimFactory.CreateClaim(user, false))
	if err != nil {
		s.logger.Errorf("UserLogin handle fail, generate token err, %v", err)
		return dto.NewApiResponse[*DTO.UserLoginResponse](dto.ErrServerError, nil)
	}
	refreshToken, err := s.claimFactory.GenerateKey(s.claimFactory.CreateClaim(user, true))
	if err != nil {
		s.logger.Errorf("UserLogin handle fail, generate refresh token err, %v", err)
		return dto.NewApiResponse[*DTO.UserLoginResponse](dto.ErrServerError, nil)
	}
	userModel := &DTO.BaseUserInfo{}
	userModel.FromUserEntity(user)

	return dto.NewApiResponse[*DTO.UserLoginResponse](
		dto.SuccessHandleRequest,
		&DTO.UserLoginResponse{
			User:         userModel,
			Token:        token,
			ExpiresIn:    int(s.claimFactory.GetJWTConfig().ExpireDuration / time.Second),
			RefreshToken: refreshToken,
		},
	)
}

func (s *AuthService) FsdLogin(form *DTO.UserFsdLogin) *DTO.UserFsdLoginResponse {
	userId := repository.GetUserId(form.Cid)
	user, err := userId.GetUser(s.userRepo)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			s.logger.Errorf("FsdLogin handle fail, %s user not found", form.Cid)
			return &DTO.UserFsdLoginResponse{Success: false, ErrorMsg: "username or password incorrect"}
		}
		s.logger.Errorf("FsdLogin handle fail, get user err, %v", err)
		return &DTO.UserFsdLoginResponse{Success: false, ErrorMsg: "Server error"}
	}

	updates := make(map[string]interface{})

	if user.Banned && user.BannedUntil.Valid && user.BannedUntil.Time.Before(time.Now()) {
		user.Banned = false
		user.BannedUntil.Valid = false
		updates["banned"] = user.Banned
		updates["banned_until"] = user.BannedUntil
	}

	if user.Banned {
		if user.BannedUntil.Valid {
			return &DTO.UserFsdLoginResponse{
				Success:  false,
				ErrorMsg: fmt.Sprintf("you were banned from the server, unban time: %s", user.BannedUntil.Time.Format("2006-01-02 15:04:05")),
			}
		}

		return &DTO.UserFsdLoginResponse{Success: false, ErrorMsg: "you were banned from the server"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		return &DTO.UserFsdLoginResponse{Success: false, ErrorMsg: "username or password incorrect"}
	}

	user.LastLoginTime.Valid = true
	user.LastLoginTime.Time = time.Now()
	user.LastLoginIP = &form.Ip
	updates["last_login_time"] = user.LastLoginTime
	updates["last_login_ip"] = user.LastLoginIP

	if err := s.userRepo.Update(user, updates); err != nil {
		s.logger.Errorf("FsdLogin handle fail, save user err, %v", err)
		return &DTO.UserFsdLoginResponse{Success: false, ErrorMsg: "Server error"}
	}

	token, err := s.claimFactory.GenerateKey(s.claimFactory.CreateFsdClaim(user))
	if err != nil {
		s.logger.Errorf("FsdLogin handle fail, generate token err, %v", err)
		return &DTO.UserFsdLoginResponse{Success: false, ErrorMsg: "Server error"}
	}

	return &DTO.UserFsdLoginResponse{Success: true, ErrorMsg: "", Token: token}
}

func (s *AuthService) RefreshToken(form *DTO.RefreshToken) *dto.ApiResponse[*DTO.RefreshTokenResponse] {
	user, err := s.userRepo.GetById(form.Uid)
	if err != nil {
		// 已签发JWT的用户必定存在, 此处为数据库错误
		s.logger.Errorf("RefreshToken handle fail, get user err, %v", err)
		return dto.NewApiResponse[*DTO.RefreshTokenResponse](dto.ErrServerError, nil)
	}

	updates := make(map[string]interface{})

	if user.Banned && user.BannedUntil.Valid && user.BannedUntil.Time.Before(time.Now()) {
		user.Banned = false
		user.BannedUntil.Valid = false
		updates["banned"] = user.Banned
		updates["banned_until"] = user.BannedUntil
	}

	// 如果在Token生效期中被封禁, 则拒绝刷新请求
	if user.Banned {
		if user.BannedUntil.Valid {
			return dto.NewApiResponse[*DTO.RefreshTokenResponse](
				dto.NewApiStatus(
					"USER_BANNED",
					fmt.Sprintf("您已被封禁，解封时间：%s", user.BannedUntil.Time.Format("2006-01-02 15:04:05")),
					dto.HttpCodePermissionDenied,
				),
				nil,
			)
		}

		return dto.NewApiResponse[*DTO.RefreshTokenResponse](service.ErrUserBanned, nil)
	}

	if err := s.userRepo.Update(user, updates); err != nil {
		s.logger.Errorf("RefreshToken handle fail, save user err, %v", err)
		return dto.NewApiResponse[*DTO.RefreshTokenResponse](dto.ErrServerError, nil)
	}

	token, err := s.claimFactory.GenerateKey(s.claimFactory.CreateClaim(user, false))
	if err != nil {
		s.logger.Errorf("RefreshToken handle fail, generate token err, %v", err)
		return dto.NewApiResponse[*DTO.RefreshTokenResponse](dto.ErrServerError, nil)
	}

	var refreshToken string
	if !form.Force && form.Raw.ExpiresAt.Add(-2*s.claimFactory.GetJWTConfig().ExpireDuration).After(time.Now()) {
		refreshToken = ""
	} else {
		refreshToken, err = s.claimFactory.GenerateKey(s.claimFactory.CreateClaim(user, true))
		if err != nil {
			s.logger.Errorf("RefreshToken handle fail, generate refresh token err, %v", err)
			return dto.NewApiResponse[*DTO.RefreshTokenResponse](dto.ErrServerError, nil)
		}
	}

	userModel := &DTO.BaseUserInfo{}
	userModel.FromUserEntity(user)

	return dto.NewApiResponse[*DTO.RefreshTokenResponse](
		dto.SuccessHandleRequest,
		&DTO.RefreshTokenResponse{
			User:         userModel,
			Token:        token,
			ExpiresIn:    int(s.claimFactory.GetJWTConfig().ExpireDuration / time.Second),
			RefreshToken: refreshToken,
		},
	)
}
