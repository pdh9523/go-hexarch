package adapter

import (
	"context"
	"errors"
	userPort "github.com/pdh9523/go-hexarch/domains/user/application/port/out"
	"github.com/pdh9523/go-hexarch/domains/user/domain"
	"github.com/pdh9523/go-hexarch/infrastructure/gorm/mapper"
	"github.com/pdh9523/go-hexarch/infrastructure/gorm/model"
	"gorm.io/gorm"
)

type UserRepositoryAdapter struct {
	db *gorm.DB
}

func NewUserRepositoryAdapter(db *gorm.DB) userPort.UserRepositoryPort {
	return &UserRepositoryAdapter{
		db: db,
	}
}

func (u *UserRepositoryAdapter) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var userModel model.User
	if err := u.db.WithContext(ctx).Where("id = ?", id).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ToDomainUser(&userModel)
}

func (u *UserRepositoryAdapter) FindByNickname(ctx context.Context, nickname string) (*domain.User, error) {
	var userModel model.User
	if err := u.db.WithContext(ctx).Where("nickname = ?", nickname).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ToDomainUser(&userModel)
}

func (u *UserRepositoryAdapter) ExistsByNickname(ctx context.Context, nickname string) (bool, error) {
	var count int64

	if err := u.db.WithContext(ctx).Model(&model.User{}).Where("nickname = ?", nickname).Count(&count).Error; err != nil {
		return false, mapper.ToDomainError(err)
	}
	return count > 0, nil
}

func (u *UserRepositoryAdapter) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := u.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, mapper.ToDomainError(err)
	}
	return count > 0, nil
}

func (u *UserRepositoryAdapter) Save(ctx context.Context, user *domain.User) (*domain.User, error) {
	userModel := mapper.ToModelUser(user)

	if err := u.db.WithContext(ctx).Create(&userModel).Error; err != nil {
		return nil, err
	}
	return mapper.ToDomainUser(userModel)
}

func (u *UserRepositoryAdapter) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	var existingUser model.User
	if err := u.db.WithContext(ctx).First(&existingUser, "id = ?", user.ID).Error; err != nil {
		return nil, mapper.ToDomainError(err)
	}
	existingUser.Nickname = user.Nickname.ToString()
	existingUser.Username = user.Username.ToString()

	if err := u.db.WithContext(ctx).Save(&existingUser).Error; err != nil {
		return nil, mapper.ToDomainError(err)
	}

	userDomain, err := mapper.ToDomainUser(&existingUser)
	if err != nil {
		return nil, err
	}
	return userDomain, nil
}

func (u *UserRepositoryAdapter) DeleteByID(ctx context.Context, id string) error {
	if err := u.db.WithContext(ctx).Where("id =?", id).Delete(&model.User{}).Error; err != nil {
		return err
	}
	return nil
}
