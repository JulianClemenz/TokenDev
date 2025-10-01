package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Workout struct {
	UserID    primitive.ObjectID `bson:"user_id,omitempty" json:"user_id" binding:"required"`
	RoutineID primitive.ObjectID `bson:"routine_id,omitempty" json:"routine_id" binding:"required"`
	Date      time.Time          `bson:"date_and_hours" json:"date_and_hours"`
}
