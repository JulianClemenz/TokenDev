package services

import (
	"AppFitness/dto"
	"AppFitness/repositories"
	"AppFitness/utils"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExcerciseInterface interface { //POST, PUT y DELETE son accesibles solo por admins (reciben actor)
	PostExcercise(excercise *dto.ExcerciseRegisterDTO, actor Actor) (*dto.ExcerciseResponseDTO, error)
	PutExcercise(excercise *dto.ExcerciseRegisterDTO, actor Actor) (*dto.ExcerciseResponseDTO, error)
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

func (service *ExcerciseService) PutExcercise(excerciseDto *dto.ExcerciseRegisterDTO, actor Actor) (*dto.ExcerciseResponseDTO, error) {

	//VALIDACIONES
	// Validacion de permisos de usuario
	if actor.Role != "admin" {
		return nil, fmt.Errorf("no tienes permisos para modificar un ejercicio, tu rol es: %s", actor.Role)
	}

	user, err := service.ExcerciseRepository.GetExcerciseByID(excerciseDto.ID)
	if err != nil {
		return nil, fmt.Errorf("no se pudo encontrar el ejercicio a modificar: %w", err)
	}
	if user.ID.IsZero() {
		return nil, fmt.Errorf("no existe un ejercicio con dicho ID")
	}

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

	excerciseModel := dto.GetModelExcerciseRegister(excerciseDto) //convertimos el dto a modelo para el repository
	result, err := service.ExcerciseRepository.PutExcercise(*excerciseModel)
	if err != nil {
		return nil, fmt.Errorf("no se pudo modificar el ejercicio: %w", err)
	}

}

type Actor struct { //solo para desarrollar, asi llegaran los datos del token desde middleware para comprobar permisos en cada metodo
	ID   string //DUDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
	Role string
}
