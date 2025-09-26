package models

import "time"

type Workout struct {
	UserID    int       `bson:"user_id" json:"user_id" binding:"required"`
	RoutineID int       `bson:"routine_id" json:"routine_id" binding:"required"`
	Date      time.Time `bson:"date_and_hours" json:"date_and_hours"`
}
