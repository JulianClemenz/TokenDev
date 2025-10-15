package services

import (
	"AppFitness/dto"
	"AppFitness/repositories"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExcerciseInterface interface {
	PostExcercise(excercise *dto.ExcerciseRegisterDTO, actor Actor) (*dto.ExcerciseResponseDTO, error)
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
	excerciseDto.CreatorUserID = actor.ID //asignamos el ID del usuario que crea el ejercicio

	excerciseModel := dto.GetModelExcerciseRegister(excerciseDto)             //convertimos el dto para registrar en model
	result, err := service.ExcerciseRepository.PostExcercise(*excerciseModel) //ejecutamos post en repository
	if err != nil {
		return nil, err
	}

	excerciseModel.ID = result.InsertedID.(primitive.ObjectID)        //asignamos el ID generado por MongoDB al modelo
	excerciseResponse := dto.NewExcerciseResponseDTO(*excerciseModel) //convertimos el modelo a dto para la respuesta
	return excerciseResponse, nil
}

type Actor struct { //solo para desarrollar, asi llegaran los datos del token desde middleware para comprobar permisos en cada metodo
	ID   string //DUDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
	Role string
}
