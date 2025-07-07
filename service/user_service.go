package service

import (
	"ati-study-jwt/model"
	"ati-study-jwt/repository"
	"ati-study-jwt/utils"
	"errors"
	"gorm.io/gorm"
	"time"
)

type UserService struct {
	UserRepo repository.UserRepository
}

type LoginData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) Register(username string, password string, birthday time.Time) (string, error) {
	// 检查用户名是否已存在
	_, err := s.UserRepo.GetUserByUsername(username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err // 数据库查询错误
	}

	if err == nil {
		return "", errors.New("username already exists") // 用户名已存在
	}

	// 创建新用户
	user := &model.User{
		Username: username,
		Password: password, // 直接存储明文密码（不推荐，仅按要求实现）
		BirthDay: birthday,
	}
	err = s.UserRepo.CreateUser(user)
	if err != nil {
		return "", err
	}

	// 生成 JWT Token
	token, err := utils.GenerateToken(user.Username)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (s *UserService) Login(username, password string) (string, error) {
	user, err := s.UserRepo.GetUserByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.New("user not found")
		}
		return "", err
	}

	if user.Password != password {
		return "", errors.New("invalid password")
	}

	token, err := utils.GenerateToken(user.Username)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
