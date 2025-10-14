package services

import (
	"AppFitness/dto"
	"AppFitness/repositories"
	"AppFitness/utils"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserInterface interface {
	PostUser(user *dto.UserRegisterDTO) (bool, error)
	GetUsers() []*dto.UserResponseDTO
	GetUserByID(id string) (*dto.UserResponseDTO, error)
	PutUser(user *dto.UserModifyDTO) (*dto.UserModifyResponseDTO, error)
}

type UserService struct {
	UserRepository repositories.UserRepositoryInterface
}

func NewUserService(UserRepository repositories.UserRepositoryInterface) *UserService {
	return &UserService{
		UserRepository: UserRepository,
	}
}

// REGISTRAR USUARIO
func (service *UserService) PostUser(userDto *dto.UserRegisterDTO) (*dto.UserResponseDTO, error) {

	if len(strings.TrimSpace(userDto.Password)) < 7 { //comprobamos que la contraseña tenga al menos 7 carateres
		return nil, fmt.Errorf("la contraseña deber tener 7 o más caracteres")
	}

	if userDto.Weight < 0 { //comprobamos que el peso ingresado no sea negativo
		return nil, fmt.Errorf("tu peso no puede ser menor a 0")
	}

	if userDto.BirthDate.After(time.Now()) { //comprobamos que la feha de nacimiento no sea mayor a hoy
		return nil, fmt.Errorf("error en fecha de nacimiento")
	}

	userDB := userDto.GetModelUserRegister()           //convertimos el dto para registrar en model
	hashed, err := utils.HashPassword(userDB.Password) //hasheamos la contraseña

	if err != nil { //comprobamos que no suceda ningun error en el hasheo de la contraseña
		return nil, fmt.Errorf("error al hashear contraseña: %w", err)
	}

	userDB.Password = hashed                             //hasheamos la contraseña

	//CAMBIAR ESTO
	usersExist, err := service.UserRepository.GetUsers() //traemos todo los usuarios para hacer comprobaciones de que no esten repetidos algunos campos

	if err != nil {
		return nil, fmt.Errorf("error al obtener usuarios: %w", err)
	}

	for _, user := range usersExist {
		if strings.EqualFold(strings.TrimSpace(user.UserName), strings.TrimSpace(userDto.UserName)) { //EqualFold no distingue mayúsculas/minúsculas, compara dos cadenas
			return nil, fmt.Errorf("ya existe ese user name")
		}
		if strings.EqualFold(strings.TrimSpace(user.Email), strings.TrimSpace(userDto.Email)) { //TrimSpace Quita espacios al principio y final de la cadena
			return nil, fmt.Errorf("email ya existente")
		}
	}

	//llamada al repositorio para guardar el usuario
	result, err := service.UserRepository.PostUser(userDB)
	if err != nil {
		return nil, fmt.Errorf("error al registrar usuario: %w", err)
	}

	userDB.ID = result.InsertedID.(primitive.ObjectID) //asignamos el ID generado por Mongo al user que vamos a devolver
	userResponse := dto.NewUserResponseDTO(userDB)     //convertimos a dto para devolver

	return userResponse, nil
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

func (services *UserService) GetUsersByID(id string) (*dto.UserResponseDTO, error) {
	userDB, err := services.UserRepository.GetUsersByID(id) //recuperamos user como model

	switch err {
	case nil:
		return dto.NewUserResponseDTO(userDB), nil //convertimos a dto y retornamos
	case mongo.ErrNoDocuments:
		fmt.Errorf("no se encontro el usuario") //user no encontrado
		return nil, err
	default:
		fmt.Errorf("error al obtener usuario por id: %w", err) //otro error
		return nil, err
		// otro
	}
}

func (services *UserService) PutUser(userDto *dto.UserModifyDTO) (*dto.UserModifyResponseDTO, error) { //aca necesitas encontrar el usuario por medio del ID, tenemos q agregar el atributo al UserModifyDTO, y desp localizar el ID en el handler para guardalo en ese dto
	
	//CAMBIAR ESTO
	//validamos que no haya campos repetidos
	usersExist, err := services.UserRepository.GetUsers() //traemos todo los usuarios para hacer comprobaciones de que no esten repetidos algunos campos
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuarios: %w", err)
	}

	for _, user := range usersExist {
		if strings.EqualFold(strings.TrimSpace(user.UserName), strings.TrimSpace(userDto.UserName)) { //EqualFold no distingue mayúsculas/minúsculas, compara dos cadenas
			return nil, fmt.Errorf("ya existe ese user name")
		}
		if strings.EqualFold(strings.TrimSpace(user.Email), strings.TrimSpace(userDto.Email)) { //TrimSpace Quita espacios al principio y final de la cadena
			return nil, fmt.Errorf("email ya existente")
		}
	}

	//convertimos el dto para registrar en model
	userDB := dto.GetModelUserModify(userDto)
	//llamada al repositorio para actualizar el usuario
	result, err := services.UserRepository.PutUser(*userDB)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar usuario: %w", err)
	}

	//recuperamos el id de result y lo asignamos al userDB para buscar el user modificado y devolverlo en el dto
	ObjID, _ := primitive.ObjectIDFromHex(result.ID.Hex())
	userModify := 

	return dto.NewUserModifyResponseDTO(*userDB), nil
}
