package services

import (
	"AppFitness/dto"
	"AppFitness/repositories"
	"fmt"
)

type RoutineInterface interface {
	PostRoutine(routine *dto.RoutineRegisterDTO) (*dto.RoutineResponseDTO, error)
}

type RoutineService struct {
	RoutineRepository repositories.RoutineRepositoryInterface
}

func NewRoutineService(routineRepository repositories.RoutineRepositoryInterface) *RoutineService {
	return &RoutineService{
		RoutineRepository: routineRepository,
	}
}

func (service RoutineService) PostRoutine(routineDTO *dto.RoutineRegisterDTO /*CreatorUserId se setea en handler*/) (*dto.RoutineResponseDTO, error) {
	//VALIDACIONES
	//validacion de campos obligatorios
	if routineDTO.Name == "" {
		return nil, fmt.Errorf("el nombre de la rutina no puede estar vacío")
	}
	if routineDTO.CreatorUserID == "" {
		return nil, fmt.Errorf("el ID del usuario creador no puede estar vacío")
	}

	//LOGICA
	result, err := service.RoutineRepository.PostRoutine(routineDTO.GetModelRoutineRegisterDTO()) //insertamos la rutina en la base de datos
	if err != nil {
		return nil, fmt.Errorf("error al crear la rutina en RoutineService.PostRoutine(): %v", err)
	}

	routineModel, err := service.RoutineRepository.GetRoutineByID(result.InsertedID.(string)) //obtenemos la rutina creada para devolverla en el response
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina creada en RoutineService.PostRoutine(): %v", err)
	}

	return dto.NewRoutineResponseDTO(routineModel), nil
}
