package services

import (
	"AppFitness/dto"
	"AppFitness/repositories"
	"fmt"
	"time"
)

type WorkoutInterface interface {
	PostWorkout(*dto.WorkoutRegisterDTO) (*dto.WorkoutResponseDTO, error)
}

type WorkoutService struct {
	WorkoutRepository repositories.WorkoutRepositoryInterface
	RoutineRepository repositories.RoutineRepositoryInterface
}

func NewWorkoutService(workoutRepository repositories.WorkoutRepositoryInterface) *WorkoutService {
	return &WorkoutService{
		WorkoutRepository: workoutRepository,
	}
}

func (ws WorkoutService) PostWorkout(workoutDTO *dto.WorkoutRegisterDTO /*UserID se setea en handler*/) (*dto.WorkoutResponseDTO, error) {
	//validaciones - GetRoutineByID(id string) (*models.Routine, error)
	result, err := ws.RoutineRepository.GetRoutineByID(workoutDTO.RoutineID)
	if err != nil {
		return nil, fmt.Errorf("rutina no encontrada: %w", err)
	}
	if result.ID.IsZero() {
		return nil, fmt.Errorf("rutina no encontrada")
	}

	workoutModel := dto.GetModelWorkoutRegisterDTO(workoutDTO)
	workoutModel.Date = time.Now()

	insertResult, err := ws.WorkoutRepository.PostWorkout(workoutModel)
	if err != nil {
		return nil, fmt.Errorf("error al crear el workout: %w", err)
	}
	if insertResult.InsertedID == nil {
		return nil, fmt.Errorf("no se pudo crear el workout")
	}

	//recuperar workout creado
	createdWorkout, err := ws.WorkoutRepository.GetWorkoutByID(insertResult.InsertedID.(string))
	if err != nil {
		return nil, fmt.Errorf("error al obtener el workout creado: %w", err)
	}
	if createdWorkout.ID.IsZero() {
		return nil, fmt.Errorf("workout creado no encontrado")
	}

	//convertir a dto y devolver
	workoutResponse := dto.NewWorkoutResponseDTO(createdWorkout)

	return workoutResponse, nil
}
