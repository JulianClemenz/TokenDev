package dto

import (
	"AppFitness/models"
	"AppFitness/utils"
	"time"
)

type RoutineRegisterDTO struct {
	Name          string
	CreatorUserID string
}
type ExcerciseInRoutineDTO struct {
	ExcerciseID string  `json:"exercise_id" binding:"required"`
	Repetitions int     `json:"repetitions" binding:"required,gt=0,lte=100"`
	Series      int     `json:"series" binding:"required,gt=0,lte=20"`
	Weight      float64 `json:"weight" binding:"gte=0,lte=1000"`
}

func GetModelRoutineRegisterDTO(routine *RoutineRegisterDTO) *models.Routine {
	return &models.Routine{
		Name:          routine.Name,
		CreatorUserID: utils.GetObjectIDFromStringID(routine.CreatorUserID),
	}
}
func GetModelExerciseInRoutineDTO(excercise *ExcerciseInRoutineDTO) models.ExcerciseInRoutine {
	return models.ExcerciseInRoutine{
		ExcerciseID: utils.GetObjectIDFromStringID(excercise.ExcerciseID),
		Repetitions: excercise.Repetitions,
		Series:      excercise.Series,
		Weight:      excercise.Weight,
	}
}
func newExcerciseInRoutineResponseDTO(excerciseList []models.ExcerciseInRoutine) []ExcerciseInRoutineDTO {
	var excerciseInRoutineDTOList []ExcerciseInRoutineDTO
	for _, excercise := range excerciseList {
		excerciseDTO := ExcerciseInRoutineDTO{
			ExcerciseID: utils.GetStringIDFromObjectID(excercise.ExcerciseID),
			Repetitions: excercise.Repetitions,
			Series:      excercise.Series,
			Weight:      excercise.Weight,
		}
		excerciseInRoutineDTOList = append(excerciseInRoutineDTOList, excerciseDTO)
	}
	return excerciseInRoutineDTOList
}

type RoutineResponseDTO struct {
	ID              string
	Name            string
	CreatorUserID   string
	ExcerciseList   []ExcerciseInRoutineDTO
	EditionDate     time.Time
	EliminationDate time.Time
	CreationDate    time.Time
}

func NewRoutineResponseDTO(routine models.Routine) *RoutineResponseDTO {
	return &RoutineResponseDTO{
		ID:              utils.GetStringIDFromObjectID(routine.ID),
		Name:            routine.Name,
		CreatorUserID:   utils.GetStringIDFromObjectID(routine.CreatorUserID), //check
		ExcerciseList:   newExcerciseInRoutineResponseDTO(routine.ExcerciseList),
		EditionDate:     routine.EditionDate,
		EliminationDate: routine.EliminationDate,
		CreationDate:    routine.CreationDate,
	}
}

type RoutineModifyDTO struct {
	IDRoutine string
	Name      string
}

type RoutineRemoveDTO struct {
	IDExercise string `json:"exercise_id" binding:"required"`
	IDRoutine  string `json:"reoutine_id" binding:"required"`
}

type ExcerciseInRoutineModifyDTO struct {
	RoutineID   string
	ExcerciseID string
	Repetitions int     `json:"repetitions" binding:"required,gt=0,lte=100"`
	Series      int     `json:"series" binding:"required,gt=0,lte=20"`
	Weight      float64 `json:"weight" binding:"gte=0,lte=1000"`
}

func GetModelFromExerciseInRoutineModifyDTO(excercise *ExcerciseInRoutineModifyDTO) models.ExcerciseInRoutine {
	return models.ExcerciseInRoutine{
		Repetitions: excercise.Repetitions,
		Series:      excercise.Series,
		Weight:      excercise.Weight,
	}
}
