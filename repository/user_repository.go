package repository

import (
	"ati-study-jwt/model"
	"context"

	"gorm.io/gorm"
)

type UserRepository struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) (err error) {
	err = r.DB.WithContext(ctx).Create(user).Error
	return
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *model.User) (err error) {
	err = r.DB.WithContext(ctx).Updates(user).Error
	return
}

func (r *UserRepository) DeleteUser(ctx context.Context, user *model.User) (err error) {
	err = r.DB.WithContext(ctx).Delete(user).Error
	return
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.DB.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.DB.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
