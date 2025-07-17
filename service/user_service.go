package service

import (
	"ati-study-jwt/middleware"
	"ati-study-jwt/model"
	"ati-study-jwt/repository"
	"context"
	"errors"
	"time"
)

type UserService struct {
	userRepo *repository.UserRepository
}

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Register(username, password string, birthday time.Time) (string, error) {
	// 检查用户是否已存在
	existingUser, _ := s.userRepo.GetUserByUsername(context.Background(), username)
	if existingUser != nil {
		return "", errors.New("user already exists")
	}

	// 创建新用户
	user := &model.User{
		Username: username,
		Password: password, // 注意：实际项目中应该对密码进行哈希处理
		BirthDay: birthday,
	}

	if err := s.userRepo.CreateUser(context.Background(), user); err != nil {
		return "", err
	}

	// 生成JWT token
	token, err := middleware.GenerateToken(username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) Login(username, password string) (string, error) {
	// 查找用户
	user, err := s.userRepo.GetUserByUsername(context.Background(), username)
	if err != nil {
		return "", errors.New("user not found")
	}

	// 验证密码（简单比较，实际项目中应该使用哈希验证）
	if user.Password != password {
		return "", errors.New("invalid password")
	}

	// 生成JWT token
	token, err := middleware.GenerateToken(username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) GetBirthday(username string) (string, error) {
	user, err := s.userRepo.GetUserByUsername(context.Background(), username)
	if err != nil {
		return "", errors.New("user not found")
	}
	return user.BirthDay.Format("2006-01-02"), nil
}

func (s *UserService) CheckUserExists(username string) bool {
	user, err := s.userRepo.GetUserByUsername(context.Background(), username)
	return err == nil && user != nil
}
