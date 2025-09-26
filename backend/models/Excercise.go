package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Excercise struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name            string             `bson:"name" json:"name" binding:"required"`
	Description     string             `bson:"description" json:"description" binding:"required"`
	CreatorUserID   int                `bson:"creator_user_id" json:"creator_user_id" binding:"required"` //int or primitive.ObjectID?
	Category        CategoryLevel      `bson:"category" json:"category" binding:"required, oneof=strength cardio flexibility balance"`
	MainMuscleGroup string             `bjson:"main_muscle_group" json:"main_mouscle_group" binding:"required"`
	DifficultLevel  string             `bjson:"difficult_level" json:"difficult_level" binding:"required"` //string or enum?
	Example         string             `bson:"example" json:"example" binding:"required"`                  //url of video
	Instructions    string             `bson:"instructions" json:"instructions" binding:"required"`
	EditionDate     time.Time          `bson:"edition_date" json:"edition_date"`
}

type CategoryLevel string

const (
	Strength    CategoryLevel = "strength"
	Cardio      CategoryLevel = "cardio"
	Flexibility CategoryLevel = "flexibility"
	Balance     CategoryLevel = "balance"
)
