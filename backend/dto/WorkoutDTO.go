package dto

import (
	"AppFitness/models"
	"AppFitness/utils"
	"time"
)

type WorkoutRegisterDTO struct {
	RoutineID string `json:"routine_id" binding:"required"`
	UserID    string
}
type WorkoutResponseDTO struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	RoutineID string    `json:"routine_id"`
	DoneAt    time.Time `json:"done_at"`
}

func GetModelWorkoutRegisterDTO(dto *WorkoutRegisterDTO) models.Workout {
	return models.Workout{
		RoutineID: utils.GetObjectIDFromStringID(dto.RoutineID),
		UserID:    utils.GetObjectIDFromStringID(dto.UserID),
	}
}

func NewWorkoutResponseDTO(workout models.Workout) *WorkoutResponseDTO {
	return &WorkoutResponseDTO{
		ID:        utils.GetStringIDFromObjectID(workout.ID),
		UserID:    utils.GetStringIDFromObjectID(workout.UserID),
		RoutineID: utils.GetStringIDFromObjectID(workout.RoutineID),
		DoneAt:    workout.Date,
	}
}
