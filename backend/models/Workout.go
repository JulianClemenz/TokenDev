package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Workout struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    int       `bson:"user_id" json:"user_id" binding:"required"`
	RoutineID int       `bson:"routine_id" json:"routine_id" binding:"required"`
	Date      time.Time `bson:"date_and_hours" json:"date_and_hours"`
}