package services

import (
	"AppFitness/dto"
	"AppFitness/repositories"
	"AppFitness/services"
	"AppFitness/utils"
)

type UserInterface interface {
	PostUser(user *dto.UserRegisterDTO) bool
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

func (service *UserService) PostUser(dto *dto.UserRegisterDTO) bool,err {
	userDB := dto.GetModelUserRegister()

	slice := strings.Split(dto.Password,"")
	if(len(slice)<)



	userDB.Password, err = utils.HashPassword(userDB.Password)

	usersExist := services.UserRepository.GetUsers()

	for _, user := range usersExist{
		if(user.UserName == dto.UserName){
			return err, fmt.Errorf("ya existe ese user name")
		}
		if(user.email == dto.email){
			return err, fmt.Errorf("email ya existente")
		}
	} 








	service.UserRepository.PostUser(userDB)


	return true, err
}

func (services *UserService) GetUsers() []*dto.UserResponseDTO {
	usersDB, _ := services.UserRepository.GetUsers()

	var users []*dto.UserResponseDTO
	for _, userDB := range usersDB{
		user := dto.NewUserResponseDTO(userDB)
		users = append(users, user)
  	}

	return users
}

func (services *UserService) GetUSersByID(id string) *dto.UserResponseDTO{
	userDB, err := services.UserRepository.GetUSersByID(id)

	var user *dto.UserResponseDTO
	if(err == nil) {
		user = dto.NewUserResponseDTO(userDB)
	}

	return user
}

func (services *UserService) PutUser(dto *dto.UserModifyDTO) bool{
	services.UserRepository.PostUser(dto.GetModelUserModify())
	return true
}

