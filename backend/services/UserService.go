package services

import (
	"AppFitness/repositories"
)

type UserInterface interface {
	PostUser(user *UserRegisterDTO) bool
}

type UserService struct {
	UserRepository repositories.UserRepositoryInterface
}

func NewUserService(UserRepository repositories.UserRepositoryInterface) *UserService {
	return &UserService{
		UserRepository: UserRepository,
	}
}

func (user UserService) PostUser(userDTO *UserRegisterDTO) bool {

}
