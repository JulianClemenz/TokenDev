package dto

import (
	"AppFitness/utils"
	"time"
)

/*

type Routine struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name            string               `bson:"name" json:"name" binding:"required"`
	CreatorUserID   int                  `bson:"creator_user_id" json:"creator_user_id" binding:"required"` //int or primitive.ObjectID?
	ExcerciseList   []ExcerciseInRoutine `bson:"excercise_list" json:"excercise_list" binding:"required,dive,required"`
	EditionDate     time.Time            `bson:"edition_date" json:"edition_date"`
	EliminationDate time.Time            `bson:"elimination_date" json:"elimination_date"`
	CreationDate    time.Time            `bson:"creation_date" json:"creation_date"`
}

type ExcerciseInRoutine struct {
	ExcerciseID     primitive.ObjectID `bson:"excercise_id,omitempty" json:"excercise_id" binding:"required"`
	Repetitions     int                `bson:"repetitions"  json:"repetitions"  binding:"required,min=1"`
	Series          int                `bson:"series"       json:"series"       binding:"required,min=1"`
	Weight          float64            `bson:"weight"       json:"weight"       binding:"gte=0"`
	EliminationDate time.Time          `bson:"elimination_date" json:"elimination_date"`
	CreationDate    time.Time          `bson:"creation_date" json:"creation_date"`
}

*/

type RoutineRegisterDTO struct {
	Name          string
	CreatorUserID string //a setar en handler
	ExcerciseList []ExcerciseInRoutineDTO
}
type ExcerciseInRoutineDTO struct {
	ExcerciseID string
	Repetitions int
	Series      int
	Weight      float64
}

func GetModelRoutineRegisterDTO(routine *RoutineRegisterDTO) *models.Routine {
	return &models.Routine{
		Name:          routine.Name,
		CreatorUserID: utils.GetObjectIDFromStringID(routine.CreatorUserID),
		ExcerciseList: getModelExcerciseInRoutineList(routine.ExcerciseList),
	}
}

type RoutineResponseDTO struct {
	ID              string
	Name            string
	CreatorUserID   int
	ExcerciseList   []ExcerciseInRoutineResponseDTO
	EditionDate     time.Time
	EliminationDate time.Time
	CreationDate    time.Time
}

func NewRoutineResponseDTO(routine models.Routine) *RoutineResponseDTO {
	return &RoutineResponseDTO{
		ID:              utils.GetStringIDFromObjectID(routine.ID),
		Name:            routine.Name,
		CreatorUserID:   int(routine.CreatorUserID.Hex()), //check
		ExcerciseList:   newExcerciseInRoutineResponseDTOList(routine.ExcerciseList),
		EditionDate:     routine.EditionDate,
		EliminationDate: routine.EliminationDate,
		CreationDate:    routine.CreationDate,
	}
}

type RoutineModifyDTO struct{}

type RoutineModifyResponseDTO struct{}
