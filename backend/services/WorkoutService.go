package services

import (
	"AppFitness/dto"
	"AppFitness/repositories"
	"fmt"
	"time"
)

type WorkoutInterface interface {
	PostWorkout(*dto.WorkoutRegisterDTO) (*dto.WorkoutResponseDTO, error)
	GetWorkouts(*dto.WorkoutRegisterDTO) ([]*dto.WorkoutResponseDTO, error)
	GetWorkoutByID(id string) (*dto.WorkoutResponseDTO, error)
}

type WorkoutService struct {
	WorkoutRepository repositories.WorkoutRepositoryInterface
	RoutineRepository repositories.RoutineRepositoryInterface
	UserRepository    repositories.UserRepositoryInterface
}

func NewWorkoutService(workoutRepository repositories.WorkoutRepositoryInterface) *WorkoutService {
	return &WorkoutService{
		WorkoutRepository: workoutRepository,
	}
}

func (ws WorkoutService) PostWorkout(workoutDTO dto.WorkoutRegisterDTO /*UserID se setea en handler*/) (*dto.WorkoutResponseDTO, error) {
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

// GetWorkouts obtiene todos los workouts de un usuario específico
func (ws WorkoutService) GetWorkouts(workoutDTO dto.WorkoutRegisterDTO /*UserID se setea en handler*/) ([]*dto.WorkoutResponseDTO, error) {

	//validar existencia de user
	user, err := ws.UserRepository.GetUsersByID(workoutDTO.UserID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario: %w", err)
	}
	if user.ID.IsZero() {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	//obtener workouts del user
	workoutsModel, err := ws.WorkoutRepository.GetWorkoutsByUserID(workoutDTO.UserID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener workouts: %w", err)
	}
	if len(workoutsModel) == 0 {
		return nil, fmt.Errorf("no se encontraron workouts para el usuario")
	}

	var workoutsDTO []*dto.WorkoutResponseDTO
	for _, workout := range workoutsModel {
		workoutDTO := dto.NewWorkoutResponseDTO(workout)
		workoutsDTO = append(workoutsDTO, workoutDTO)
	}

	return workoutsDTO, nil
}

// GetWorkoutByID obtiene un workout por su ID
func (ws WorkoutService) GetWorkoutByID(id string) (*dto.WorkoutResponseDTO, error) {

	//validacion de existencia de user
	user, err := ws.UserRepository.GetUsersByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario: %w", err)
	}
	if user.ID.IsZero() {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	//obtener workout por id

}
