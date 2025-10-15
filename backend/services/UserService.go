package services

import (
	"AppFitness/dto"
	"AppFitness/repositories"
	"AppFitness/utils"
	"fmt"
	"strings"
	"time"
)

type UserInterface interface {
	PostUser(user *dto.UserRegisterDTO) (bool, error)
	GetUsers() []*dto.UserResponseDTO
	GetUserByID(id string) *dto.UserResponseDTO
	PutUser(user *dto.UserModifyDTO) bool
}

type UserService struct {
	UserRepository repositories.UserRepositoryInterface
}

func NewUserService(UserRepository repositories.UserRepositoryInterface) *UserService {
	return &UserService{
		UserRepository: UserRepository,
	}
}

func (service *UserService) PostUser(dto *dto.UserRegisterDTO) (bool, error) {

	if dto.Weight < 0 { //comprobamos que el peso ingresado no sea negativo
		return false, fmt.Errorf("tu peso no puede ser menor a 0")
	}

	if dto.BirthDate.After(time.Now()) { //comprobamos que la feha de nacimiento no sea mayor a hoy
		return false, fmt.Errorf("error en fecha de nacimiento")
	}

	dto.Email = strings.ToLower(strings.TrimSpace(dto.Email))
	dto.UserName = strings.TrimSpace(dto.UserName)

	verificacionEmail, err := service.UserRepository.ExistByEmail(dto.Email)
	if err != nil {
		return false, fmt.Errorf("no se pudo verificar email: %w", err)
	}

	if verificacionEmail {
		return false, fmt.Errorf("dicho email ya existe")
	}

	verificacionUserName, err := service.UserRepository.ExistByUserName(dto.UserName)

	if err != nil {
		return false, fmt.Errorf("no se pudo verificar nombre de usuario: %w", err)
	}

	if verificacionUserName {
		return false, fmt.Errorf("dicho nombre de usuario ya existe")
	}

	userDB := dto.GetModelUserRegister() //convertimos el dto para registrar en model

	hashed, err := utils.HashPassword(userDB.Password)
	if err != nil { //comprobamos que no suceda ningun error en el hasheo de la contraseña
		return false, fmt.Errorf("error al hashear contraseña: %w", err)
	}

	userDB.Password = hashed //hasheamos la contraseña

	if _, err := service.UserRepository.PostUser(userDB); err != nil {
		return false, fmt.Errorf("error al insertar usuario: %w", err)
	}
	return true, nil
}

func (services *UserService) GetUsers() []*dto.UserResponseDTO {
	usersDB, _ := services.UserRepository.GetUsers()

	var users []*dto.UserResponseDTO
	for _, userDB := range usersDB {
		user := dto.NewUserResponseDTO(userDB)
		users = append(users, user)
	}

	return users
}

func (services *UserService) GetUSersByID(id string) *dto.UserResponseDTO {
	userDB, err := services.UserRepository.GetUsersByID(id)

	var user *dto.UserResponseDTO
	if err == nil {
		user = dto.NewUserResponseDTO(userDB)
	}

	return user
}
