package services

import (
	"AppFitness/dto"
	"AppFitness/repositories"
	"AppFitness/services"
	"AppFitness/utils"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoutineInterface interface {
	PostRoutine(routineDTO *dto.RoutineRegisterDTO) (*dto.RoutineResponseDTO, error)
	GetRoutines()([]*dto.RoutineResponseDTO, error)

}

type RoutineService struct {
	RoutineRepository repositories.RoutineRepositoryInterface
}

func NewRoutineService(routineRepository repositories.RoutineRepositoryInterface) *RoutineService {
	return &RoutineService{
		RoutineRepository: routineRepository,
	}
}

func (service *RoutineService) PostRoutine(routineDTO *dto.RoutineRegisterDTO) (*dto.RoutineResponseDTO, error) {
	routineDTO.Name = strings.ToLower(strings.TrimSpace(routineDTO.Name))

	//validacion de campos obligatorios
	if routineDTO.Name == "" {
		return nil, fmt.Errorf("el nombre de la rutina no puede estar vacío")
	}
	if routineDTO.CreatorUserID == "" {
		return nil, fmt.Errorf("el ID del usuario creador no puede estar vacío")
	}

	verificacion, err := service.RoutineRepository.ExistByRutineName(routineDTO.Name)
	if err != nil {
		return nil, fmt.Errorf("no se pudo verificar si existe una rutina con el mismo nombre: %w", err)
	}

	if verificacion {
		return nil, fmt.Errorf("dicho nombre de rutina ya existe")
	}

	//LOGICA
	model := dto.GetModelRoutineRegisterDTO(routineDTO)
	result, err := service.RoutineRepository.PostRoutine(model) //insertamos la rutina en la base de datos
	if err != nil {
		return nil, fmt.Errorf("error al crear la rutina en RoutineService.PostRoutine(): %v", err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("no se pudo obtener el ObjectID insertado")
	}
	idStr := utils.GetStringIDFromObjectID(oid)

	routineModel, err := service.RoutineRepository.GetRoutineByID(idStr)//obtenemos la rutina de tipo response creada para devolver

	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina creada en RoutineService.PostRoutine(): %v", err)
	}

	return dto.NewRoutineResponseDTO(routineModel), nil
}

func (service *RoutineService) GetRoutines()([]*dto.RoutineResponseDTO, error){
	routinesDB, err := services.RoutineRepository.GetRoutines()
	if(err != nil){
		return nil, fmt.Errorf("error al obtener rutinas %v:", err)
	}
	var routines []*dto.RoutineResponseDTO
	for _, routineDB := range routinesDB {
		routine := dto.NewRoutineResponseDTO(routineDB)
		routines = append(routines, routine)
	}

	return routines, nil
}

func (service *RoutineService) GetR
