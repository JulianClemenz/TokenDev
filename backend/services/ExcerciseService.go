package services

import (
	"AppFitness/dto"
	"AppFitness/repositories"
	"AppFitness/utils"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExcerciseInterface interface { //POST, PUT y DELETE son accesibles solo por admins (reciben actor)
	//admins
	PostExcercise(excercise *dto.ExcerciseRegisterDTO, actor Actor) (*dto.ExcerciseResponseDTO, error)
	PutExcercise(excercise *dto.ExcerciseRegisterDTO, actor Actor) (*dto.ExcerciseResponseDTO, error)
	DeleteExcercise(id string, actor Actor) (bool, error)
	//todos los usuarios
	GetExcercises() ([]*dto.ExcerciseResponseDTO, error)
	GetExcerciseByID(id string) (*dto.ExcerciseResponseDTO, error)
}

type ExcerciseService struct {
	ExcerciseRepository repositories.ExcerciseRepositoryInterface
}

func NewExcerciseService(ExcerciseRepository repositories.ExcerciseRepositoryInterface) *ExcerciseService {
	return &ExcerciseService{
		ExcerciseRepository: ExcerciseRepository,
	}
}

func (service *ExcerciseService) PostExcercise(excerciseDto *dto.ExcerciseRegisterDTO, actor Actor) (*dto.ExcerciseResponseDTO, error) {

	//VALIDACIONES
	// Validacion de permisos de usuario
	if actor.Role != "admin" {
		return nil, fmt.Errorf("no tienes permisos para crear un ejercicio, tu rol es: %s", actor.Role)
	}

	// Validaciones de campos obligatorios
	if excerciseDto.Name == "" {
		return nil, fmt.Errorf("el nombre del ejercicio no puede estar vacío")
	}
	if excerciseDto.DifficultLevel == "" {
		return nil, fmt.Errorf("el nivel de dificultad no puede estar vacío")
	}
	if excerciseDto.MainMuscleGroup == "" {
		return nil, fmt.Errorf("el grupo muscular no puede estar vacío")
	}
	if excerciseDto.Description == "" {
		return nil, fmt.Errorf("la descripción del ejercicio no puede estar vacía")
	}
	if excerciseDto.Category == "" {
		return nil, fmt.Errorf("la categoría del ejercicio no puede estar vacía")
	}

	// Validacion de existencia por nombre
	nameExist, err := service.ExcerciseRepository.ExistByName(excerciseDto.Name)
	if err != nil {
		return nil, fmt.Errorf("no se pudo verificar el nombre del ejercicio en la base de datos: %w", err)
	}
	if nameExist == true {
		return nil, fmt.Errorf("ya existe un ejercicio con ese nombre")
	}

	//LOGICA
	excerciseModel := dto.GetModelExcerciseRegister(excerciseDto)             //convertimos el dto a modelo para el repository
	excerciseModel.CreatorUserID = utils.GetObjectIDFromStringID(actor.ID)    //asignamos el ObjectID del usuario que crea el ejercicio
	result, err := service.ExcerciseRepository.PostExcercise(*excerciseModel) //ejecutamos post en repository
	if err != nil {
		return nil, err
	}

	excerciseModel.ID = result.InsertedID.(primitive.ObjectID)        //asignamos el ID generado por MongoDB al model
	excerciseResponse := dto.NewExcerciseResponseDTO(*excerciseModel) //convertimos el modelo a dto para la respuesta
	return excerciseResponse, nil
}

func (service *ExcerciseService) PutExcercise(newData *dto.ExcerciseModifyDTO, id string, actor Actor) (*dto.ExcerciseModifyResponseDTO, error) {

	//VALIDACIONES
	// Validacion de permisos de usuario
	if actor.Role != "admin" {
		return nil, fmt.Errorf("no tienes permisos para modificar un ejercicio, tu rol es: %s", actor.Role)
	}

	// Validaciones de campos obligatorios
	if newData.Name == "" {
		return nil, fmt.Errorf("datos vacios")
	}
	if id == "" {
		return nil, fmt.Errorf("el id del ejercicio no puede estar vacío")
	}

	//LOGICA
	_, err := service.ExcerciseRepository.GetExcerciseByID(id) //comprobamos que el ejercicio a modificar existe
	if err != nil {
		return nil, fmt.Errorf("error al obtener el ejercicio a modificar: %w", err)
	}

	excerciseModel := dto.GetModelExcerciseModify(newData) //convertimos el dto a modelo para el repository
	excerciseModel.ID = utils.GetObjectIDFromStringID(id)  //asignamos el ObjectID del ejercicio a modificar
	excerciseModel.EditionDate = time.Now()                //actualizamos la fecha de edición

	result, err := service.ExcerciseRepository.PutExcercise(*excerciseModel) //ejecutamos put en repository
	if err != nil {
		return nil, err
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no se modificó ningún ejercicio")
	}

	excerciseModify, err := service.ExcerciseRepository.GetExcerciseByID(id) //obtenemos el ejercicio modificado para devolverlo
	if err != nil {
		return nil, fmt.Errorf("error al obtener el ejercicio modificado: %w", err)
	}

	return dto.NewExcerciseModifyResponseDTO(excerciseModify), nil
}

func (service *ExcerciseService) GetExcercises() ([]*dto.ExcerciseResponseDTO, error) {
	excercisesDB, err := service.ExcerciseRepository.GetExcercises()
	if err != nil {
		return nil, fmt.Errorf("error al obtener ejercicios: %w", err)
	}

	var excercises []*dto.ExcerciseResponseDTO
	for _, excerciseDB := range excercisesDB {
		excercise := dto.NewExcerciseResponseDTO(excerciseDB)
		excercises = append(excercises, excercise)
	}
	return excercises, nil
}

func (service *ExcerciseService) GetExcerciseByID(id string) (*dto.ExcerciseResponseDTO, error) {
	userDB, err := service.ExcerciseRepository.GetExcerciseByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener ejercicio: %w", err)
	}
	return dto.NewExcerciseResponseDTO(userDB), nil
}

/*func (service *ExcerciseService) DeleteExcercise(id string, actor Actor) (bool, error) {}
ELIMINACION LOGICA O VERDADERA?*/

type Actor struct { //solo para desarrollar, asi llegaran los datos del token desde middleware para comprobar permisos en cada metodo
	ID   string //DUDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
	Role string
}
