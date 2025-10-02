package dto

import (
	"AppFitness/models"
	"time"
)

type UserRegisterDTO struct {
	Name       string
	LastName   string
	UserName   string
	Email      string
	Password   string
	BirthDate  time.Time
	Role       string //
	Weight     float32
	Height     float32
	Experience string //podria ser string
	Objetive   string //podria ser string
}

func NewUserRegisterDTO(user models.User) *UserRegisterDTO {
	return &UserRegisterDTO{
		Name:       user.Name,
		LastName:   user.LastName,
		UserName:   user.UserName,
		Email:      user.Email,
		Password:   user.Password,
		BirthDate:  user.BirthDate,
		Role:       string(user.Role),
		Weight:     user.Weight,
		Height:     user.Height,
		Experience: string(user.Experience),
		Objetive:   string(user.Objetive),
	}
}

func (user UserRegisterDTO) GetModel() models.User {
	return models.User{
		Name:       user.Name,
		LastName:   user.LastName,
		UserName:   user.UserName,
		Email:      user.Email,
		Password:   user.Password,
		BirthDate:  user.BirthDate,
		Role:       models.AdminLevel(user.Role),
		Weight:     user.Weight,
		Height:     user.Height,
		Experience: models.ExperienceLevel(user.Experience),
		Objetive:   models.ObjetiveLevel(user.Objetive),
	}
}
