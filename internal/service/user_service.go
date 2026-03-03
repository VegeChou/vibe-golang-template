package service

import (
	"errors"
	"strings"

	"vibe-golang-template/internal/model"
)

var ErrInvalidUserInput = errors.New("invalid user input")

type UserRepository interface {
	List() []model.User
	Create(user model.User) model.User
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) ListUsers() []model.User {
	return s.repo.List()
}

func (s *UserService) CreateUser(input model.User) (model.User, error) {
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Email) == "" {
		return model.User{}, ErrInvalidUserInput
	}
	return s.repo.Create(input), nil
}
