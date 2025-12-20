// Package repository
package repository

import (
	"errors"
	"time"
	repos "user-service/src/interfaces/repository"

	"gorm.io/gorm"
	"half-nothing.cn/service-core/database"
	"half-nothing.cn/service-core/interfaces/database/entity"
	"half-nothing.cn/service-core/interfaces/database/repository"
	"half-nothing.cn/service-core/interfaces/logger"
)

type UserRepository struct {
	*database.BaseRepository[*entity.User]
	pageReq database.PageableInterface[*entity.User]
}

func NewUserRepository(
	lg logger.Interface,
	db *gorm.DB,
	queryTimeout time.Duration,
) *UserRepository {
	return &UserRepository{
		BaseRepository: database.NewBaseRepository[*entity.User](lg, "user-repository", db, queryTimeout),
		pageReq:        database.NewPageRequest[*entity.User](db),
	}
}

func (repo *UserRepository) GetById(id uint) (*entity.User, error) {
	if id <= 0 {
		return nil, repository.ErrArgument
	}
	user := &entity.User{ID: id}
	err := repo.Query(func(tx *gorm.DB) error {
		return tx.Preload("CurrentAvatar").First(user).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = repository.ErrRecordNotFound
	}
	return user, err
}

func (repo *UserRepository) GetByIdOrCid(id uint) (*entity.User, error) {
	if id < 0 {
		return nil, repository.ErrArgument
	}
	user := &entity.User{}
	err := repo.Query(func(tx *gorm.DB) error {
		return tx.Preload("CurrentAvatar").Where("id = ? OR cid = ?", id, id).First(user).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = repos.ErrUserNotFound
	}
	return user, err
}

func (repo *UserRepository) GetByUsernameOrEmail(usernameOrEmail string) (*entity.User, error) {
	if usernameOrEmail == "" {
		return nil, repository.ErrArgument
	}
	user := &entity.User{}
	err := repo.Query(func(tx *gorm.DB) error {
		return tx.Preload("CurrentAvatar").Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(user).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = repos.ErrUserNotFound
	}
	return user, err
}
