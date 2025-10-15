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

func (service *UserService) PostUser(userDto *dto.UserRegisterDTO) (*dto.UserResponseDTO, error) {

	//VALIDACIONES
	if userDto.Weight < 0 { //comprobamos que el peso ingresado no sea negativo
		return nil, fmt.Errorf("tu peso no puede ser menor a 0")
	}

	if userDto.BirthDate.After(time.Now()) { //comprobamos que la feha de nacimiento no sea mayor a hoy
		return nil, fmt.Errorf("error en fecha de nacimiento")
	}

	userDto.Email = strings.ToLower(strings.TrimSpace(userDto.Email))
	userDto.UserName = strings.TrimSpace(userDto.UserName)

	verificacionEmail, err := service.UserRepository.ExistByEmail(userDto.Email) //verifica si existe el email en la base de datos
	if err != nil {
		return nil, fmt.Errorf("no se pudo verificar email: %w", err)
	}

	if verificacionEmail {
		return nil, fmt.Errorf("dicho email ya existe")
	}

	verificacionUserName, err := service.UserRepository.ExistByUserName(userDto.UserName) //verifica si existe el username en la base de datos

	if err != nil {
		return nil, fmt.Errorf("no se pudo verificar nombre de usuario: %w", err)
	}

	if verificacionUserName {
		return nil, fmt.Errorf("dicho nombre de usuario ya existe")
	}

	//LOGICA
	userDB := userDto.GetModelUserRegister() //convertimos el dto para registrar en model

	hashed, err := utils.HashPassword(userDB.Password)
	if err != nil { //comprobamos que no suceda ningun error en el hasheo de la contraseña
		return nil, fmt.Errorf("error al hashear contraseña: %w", err)
	}

	userDB.Password = hashed //hasheamos la contraseña

	result, err := service.UserRepository.PostUser(userDB)
	if err != nil {
		return nil, fmt.Errorf("error al insertar usuario: %w", err)
	}

	userDB.ID = result.InsertedID.(primitive.ObjectID) //asignamos el ID generado por Mongo al userDB
	userResponse := dto.NewUserResponseDTO(userDB)     //convertimos el model a dto para devolverlo
	return userResponse, nil
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

func (s *UserService) PutUser(newData *dto.UserModifyDTO) (*dto.UserModifyResponseDTO, error) {

	if newData == nil {
		return nil, fmt.Errorf("datos vacíos")
	}
	user, err := s.UserRepository.GetUsersByID(newData.ID)

	if err != nil {
		return nil, fmt.Errorf("error al buscar usuario")
	}
	if user.ID.IsZero() {
		return nil, fmt.Errorf("no existe ningun usuario con ese ID") //verificamos que existe el usuario antes de modificar
	}

	if newData.Weight < 0 { //comprobamos que el peso ingresado no sea negativo
		return nil, fmt.Errorf("tu peso no puede ser menor a 0")
	}

	newData.UserName = strings.TrimSpace(newData.UserName)

	if newData.UserName != strings.TrimSpace(user.UserName) {//solo verificamos si el user name es diferente al q ya estaba
		exist, err := s.UserRepository.ExistByUserNameExceptID(newData.ID, newData.UserName)
		if err != nil {
			return nil, fmt.Errorf("no se pudo verificar nombre de usuario: %w", err)
		}

		if exist {
			return nil, fmt.Errorf("dicho nombre de usuario ya existe")
		}
	}

	userDB := dto.GetModelUserModify(newData)
	userResp := dto.NewUserModifyResponseDTO(userDB)

	if _, err := s.UserRepository.PutUser(userDB); err != nil {
		return nil, fmt.Errorf("error al modificar usuario: %w", err)
	}

	return userResp, nil
}
