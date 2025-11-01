package dto

import (
	"AppFitness/models"
	"AppFitness/utils"
	"time"
)

type WorkoutRegisterDTO struct {
	RoutineID   string `json:"routine_id" binding:"required"`
	RoutineName string
	UserID      string
}
type WorkoutResponseDTO struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	RoutineID   string    `json:"routine_id"`
	RoutineName string    `json: routine_name`
	DoneAt      time.Time `json:"done_at"`
}

func GetModelWorkoutRegisterDTO(dto *WorkoutRegisterDTO) models.Workout {
	return models.Workout{
		RoutineID:   utils.GetObjectIDFromStringID(dto.RoutineID),
		UserID:      utils.GetObjectIDFromStringID(dto.UserID),
		RoutineName: dto.RoutineName,
	}
}

func NewWorkoutResponseDTO(workout models.Workout) *WorkoutResponseDTO {
	return &WorkoutResponseDTO{
		ID:          utils.GetStringIDFromObjectID(workout.ID),
		UserID:      utils.GetStringIDFromObjectID(workout.UserID),
		RoutineID:   utils.GetStringIDFromObjectID(workout.RoutineID),
		RoutineName: workout.RoutineName,
		DoneAt:      workout.Date,
	}
}

type WorkoutStatsDTO struct {
	TotalWorkouts    int                //cantidad total de workouts del user
	WeeklyFrequency  float64            // promedio de entrenamientos desde que se realizop el primero (ir contando la cantidad de dias que hay entre entrenamientos (desde el primero hasta el ult) y dividir por la cantidad de entrenamientos)
	MostUsedRoutines []RoutineUsageDTO  //ranking de rutinas mas usadas
	ProgressOverTime []ProgressPointDTO //para grafica entrenamientos-dias
}

type RoutineUsageDTO struct {
	RoutineName string
	Count       int
}

type ProgressPointDTO struct {
	Date  string
	Count int
}

type WorkoutDeleteDTO struct {
	RoutineID string
	UserID    string
}
