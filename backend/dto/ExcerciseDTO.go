package dto

import (
	"AppFitness/models"
	"AppFitness/utils"
	"time"
)

type ExcerciseRegisterDTO struct {
	Name            string
	Description     string
	CreatorUserID   string
	Category        string
	MainMuscleGroup string
	DifficultLevel  string
	Example         string
	Instructions    string
}

func GetModelExcerciseRegister(excercise *ExcerciseRegisterDTO) *models.Excercise {
	return &models.Excercise{
		Name:            excercise.Name,
		Description:     excercise.Description,
		CreatorUserID:   utils.GetObjectIDFromStringID(excercise.CreatorUserID),
		Category:        models.CategoryLevel(excercise.Category),
		MainMuscleGroup: excercise.MainMuscleGroup,
		DifficultLevel:  excercise.DifficultLevel,
		Example:         excercise.Example,
		Instructions:    excercise.Instructions,
	}
}

type ExcerciseResponseDTO struct {
	ID              string
	Name            string
	Description     string
	CreatorUserID   string
	Category        string
	MainMuscleGroup string
	DifficultLevel  string
	Example         string
	Instructions    string
	EditionDate     time.Time
	EliminationDate time.Time
	CreationDate    time.Time
}
