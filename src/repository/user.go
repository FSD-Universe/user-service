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
		err = repos.ErrUserNotFound
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

func (repo *UserRepository) CheckCidUsernameAndEmail(cid uint, username string, email string) (bool, error) {
	if cid <= 0 && username == "" && email == "" {
		return false, repository.ErrArgument
	}
	var count int64
	err := repo.Query(func(tx *gorm.DB) error {
		return tx.Model(&entity.User{}).
			Where("cid = ? OR username = ? OR email = ?", cid, username, email).
			Count(&count).Error
	})
	if err != nil {
		return false, err
	}
	return count == 0, err
}

func (repo *UserRepository) GetPages(pageNum int, pageSize int, search string) (users []*entity.User, total int64, err error) {
	users = make([]*entity.User, 0, pageSize)
	var queryFunc func(tx *gorm.DB) *gorm.DB
	if search != "" {
		queryFunc = func(tx *gorm.DB) *gorm.DB {
			return tx.Where("username LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%").
				Preload("CurrentAvatar").
				Order("cid")
		}
	} else {
		queryFunc = func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("CurrentAvatar").Order("cid")
		}
	}
	total, err = repo.QueryWithPagination(repo.pageReq, database.NewPage[*entity.User](pageNum, pageSize, &users, &entity.User{}, queryFunc))
	return
}
