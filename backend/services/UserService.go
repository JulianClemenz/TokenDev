package services

import (
	"AppFitness/dto"
	"AppFitness/repositories"
	"AppFitness/services"
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

	if len(strings.TrimSpace(dto.Password)) < 7 { //comprobamos que la contraseña tenga al menos 7 carateres

		return false, fmt.Errorf("la contraseña deber tener 7 o más caracteres")
	}

	if dto.Weight < 0 { //comprobamos que el peso ingresado no sea negativo
		return false, fmt.Errorf("tu peso no puede ser menor a 0")
	}

	if dto.BirthDate.After(time.Now()) { //comprobamos que la feha de nacimiento no sea mayor a hoy
		return false, fmt.Errorf("error en fecha de nacimiento")
	}

	userDB := dto.GetModelUserRegister() //convertimos el dto para registrar en model
	hashed, err := utils.HashPassword(userDB.Password)

	if err != nil { //comprobamos que no suceda ningun error en el hasheo de la contraseña
		return false, fmt.Errorf("error al hashear contraseña: %w", err)
	}

	userDB.Password = hashed                              //hasheamos la contraseña
	usersExist, err := services.UserRepository.GetUsers() //traemos todo los usuarios para hacer comprobaciones de que no esten repetidos algunos campos

	if err != nil {
		return false, fmt.Errorf("error al obtener usuarios: %w", err)
	}

	for _, user := range usersExist {
		if strings.EqualFold(strings.TrimSpace(user.UserName), strings.TrimSpace(dto.UserName)) { //EqualFold no distingue mayúsculas/minúsculas, compara dos cadenas
			return false, fmt.Errorf("ya existe ese user name")
		}
		if strings.EqualFold(strings.TrimSpace(user.Email), strings.TrimSpace(dto.Email)) { //TrimSpace Quita espacios al principio y final de la cadena
			return false, fmt.Errorf("email ya existente")
		}
	}
	service.UserRepository.PostUser(userDB)
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
	userDB, err := services.UserRepository.GetUSersByID(id)

	var user *dto.UserResponseDTO
	if err == nil {
		user = dto.NewUserResponseDTO(userDB)
	}

	return user
}

func (services *UserService) PutUser(dto *dto.UserModifyDTO) bool {
	services.UserRepository.PostUser(dto.GetModelUserModify())
	return true
}
