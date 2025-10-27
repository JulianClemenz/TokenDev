package dto

import (
	"AppFitness/models"
	"AppFitness/utils"
	"time"
)

type RoutineRegisterDTO struct {
	Name          string
	CreatorUserID string //a setar en handler
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
	}
}
func GetModelExerciseInRoutineDTO(excercise ExcerciseInRoutineDTO) models.ExcerciseInRoutine {
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

type RoutineModifyDTO struct{}

type RoutineModifyResponseDTO struct{}
