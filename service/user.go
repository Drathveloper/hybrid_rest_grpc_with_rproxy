package service

import (
	"errors"
	"poc_hybrid_grpc_rest/model"
)

var ErrInvalidUser = errors.New("invalid request user")

type UserService interface {
	CreateUser(user model.User) (string, error)
}

type userService struct {
}

func NewUserService() UserService {
	return &userService{}
}

func (s *userService) CreateUser(user model.User) (string, error) {
	if user.Username == "" || user.Password == "" {
		return "", ErrInvalidUser
	}
	return user.ToID(), nil
}
