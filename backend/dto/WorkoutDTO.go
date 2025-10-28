package dto

import "time"

type WorkoutRegisterDTO struct {
	RoutineID string `json:"routine_id" binding:"required"`
}
type WorkoutResponseDTO struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	RoutineID string    `json:"routine_id"`
	DoneAt    time.Time `json:"done_at"`
}
