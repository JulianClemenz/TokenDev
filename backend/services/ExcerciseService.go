package services

import (
	"AppFitness/dto"
	"AppFitness/repositories"
	"fmt"
	"strings"
)

type ExcerciseInterface interface {
	PostExcercise(excercise *dto.ExcerciseRegisterDTO) (*dto.ExcerciseResponseDTO, error)
}

type ExcerciseService struct {
	ExcerciseRepository repositories.ExcerciseRepositoryInterface
}

func NewExcerciseService(ExcerciseRepository repositories.ExcerciseRepositoryInterface) *ExcerciseService {
	return &ExcerciseService{
		ExcerciseRepository: ExcerciseRepository,
	}
}

// REGISTRAR EJERCICIO
func (service *ExcerciseService) PostExcercise(excerciseDto *dto.ExcerciseRegisterDTO) (*dto.ExcerciseResponseDTO, error) {
	if excerciseDto.Name == "" {
		return nil, fmt.Errorf("el nombre del ejercicio no puede estar vacío")
	}
	if excerciseDto.DifficultLevel == "" {
		return nil, fmt.Errorf("el nivel de dificultad no puede estar vacío")
	}
	if excerciseDto.CreatorUserID == "" {
		return nil, fmt.Errorf("el ID del usuario creador no puede estar vacío")
	}

	excerciseModel := dto.GetModelExcerciseRegister(excerciseDto) //convertimos el dto para registrar en model

	excercisesList, err := service.ExcerciseRepository.GetExcercises() //traemos todo los ejercicios para hacer comprobaciones de que no esten repetidos algunos campos
	if err != nil {
		return nil, fmt.Errorf("error al obtener ejercicios: %w", err)
	}

	for _, excercise := range excercisesList {
		if strings.EqualFold(excercise.Name, excerciseDto.Name) { //EqualFold no distingue mayúsculas/minúsculas, compara dos cadenas
			return nil, fmt.Errorf("ya existe un ejercicio con ese nombre")
		}
	}

	/*
		for _, user := range usersExist {
			if strings.EqualFold(strings.TrimSpace(user.UserName), strings.TrimSpace(dto.UserName)) { //EqualFold no distingue mayúsculas/minúsculas, compara dos cadenas
				return false, fmt.Errorf("ya existe ese user name")
			}
			if strings.EqualFold(strings.TrimSpace(user.Email), strings.TrimSpace(dto.Email)) { //TrimSpace Quita espacios al principio y final de la cadena
				return false, fmt.Errorf("email ya existente")
			}
		}
	*/
}

type Actor struct { //solo para desarrollar, asi llegaran los datos del token desde middleware para comprobar permisos en cada metodo
	ID   string
	Role string
}
