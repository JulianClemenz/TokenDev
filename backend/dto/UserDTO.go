package dto

import (
	"AppFitness/dto"
	"AppFitness/models"
	"AppFitness/utils"
	"time"
)

type UserRegisterDTO struct {
	Name       string
	LastName   string
	UserName   string `json:"user_name" binding:"required,min=5"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=7"`
	BirthDate  time.Time
	Role       string
	Weight     float32
	Height     float32
	Experience string
	Objetive   string
}

func (user UserRegisterDTO) GetModelUserRegister() models.User {
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

type UserResponseDTO struct {
	ID         string
	Name       string
	LastName   string
	UserName   string
	Email      string
	BirthDate  time.Time
	Weight     float32
	Height     float32
	Experience string
	Objetive   string
}

func NewUserResponseDTO(user models.User) *UserResponseDTO {
	return &UserResponseDTO{
		ID:         utils.GetStringIDFromObjectID(user.ID),
		Name:       user.Name,
		LastName:   user.LastName,
		UserName:   user.UserName,
		Email:      user.Email,
		BirthDate:  user.BirthDate,
		Weight:     user.Weight,
		Height:     user.Height,
		Experience: string(user.Experience),
		Objetive:   string(user.Objetive),
	}
}

type UserModifyDTO struct {
	ID         string
	UserName   string
	Email      string `json:"email" binding:"required,email"`
	Role       string
	Weight     float32 `json:"weight" binding:"gte=0"`
	Height     float32 `json:"height" binding:"gte=0"`
	Experience string
	Objetive   string
}

func GetModelUserModify(user *dto.UserModifyDTO) models.User {
	return models.User{
		ID:         utils.GetObjectIDFromStringID(user.ID),
		UserName:   user.UserName,
		Email:      user.Email,
		Role:       models.AdminLevel(user.Role),
		Weight:     user.Weight,
		Height:     user.Height,
		Experience: models.ExperienceLevel(user.Experience),
		Objetive:   models.ObjetiveLevel(user.Objetive),
	}
}

type UserModifyResponseDTO struct {
	UserName   string
	Email      string
	Role       string
	Weight     float32
	Height     float32
	Experience string
	Objetive   string
}

func NewUserModifyResponseDTO(user models.User) *UserModifyResponseDTO {
	return &UserModifyResponseDTO{
		UserName:   user.UserName,
		Email:      user.Email,
		Role:       string(user.Role),
		Weight:     user.Weight,
		Height:     user.Height,
		Experience: string(user.Experience),
		Objetive:   string(user.Objetive),
	}
}

func GetModelUserModifyResponse(user UserModifyResponseDTO) *models.User {
	return &models.User{
		UserName:   user.UserName,
		Email:      user.Email,
		Role:       models.AdminLevel(user.Role),
		Weight:     user.Weight,
		Height:     user.Height,
		Experience: models.ExperienceLevel(user.Experience),
		Objetive:   models.ObjetiveLevel(user.Objetive),
	}
}
