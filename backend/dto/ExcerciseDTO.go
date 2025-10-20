package dto

import (
	"AppFitness/models"
	"AppFitness/utils"
	"time"
)

// ExcerciseRegisterDTO
// ExcerciseResponseDTO
// ExcerciseModifyDTO

type ExcerciseRegisterDTO struct {
	Name            string
	Description     string
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

func NewExcerciseResponseDTO(excercise models.Excercise) *ExcerciseResponseDTO {
	return &ExcerciseResponseDTO{
		ID:              utils.GetStringIDFromObjectID(excercise.ID),
		Name:            excercise.Name,
		Description:     excercise.Description,
		CreatorUserID:   utils.GetStringIDFromObjectID(excercise.CreatorUserID),
		Category:        string(excercise.Category),
		MainMuscleGroup: excercise.MainMuscleGroup,
		DifficultLevel:  excercise.DifficultLevel,
		Example:         excercise.Example,
		Instructions:    excercise.Instructions,
		EditionDate:     excercise.EditionDate,
		EliminationDate: excercise.EliminationDate,
		CreationDate:    excercise.CreationDate,
	}
}

type ExcerciseModifyDTO struct {
	ID              string
	Name            string
	Description     string
	Category        string
	MainMuscleGroup string
	DifficultLevel  string
	Example         string
	Instructions    string
}

func GetModelExcerciseModify(excercise *ExcerciseModifyDTO) *models.Excercise {
	return &models.Excercise{
		Name:            excercise.Name,
		Description:     excercise.Description,
		Category:        models.CategoryLevel(excercise.Category),
		MainMuscleGroup: excercise.MainMuscleGroup,
		DifficultLevel:  excercise.DifficultLevel,
		Example:         excercise.Example,
		Instructions:    excercise.Instructions,
	}
}

type ExcerciseModifyResponseDTO struct {
	Name            string
	Description     string
	CreatorUserID   string
	Category        string
	MainMuscleGroup string
	DifficultLevel  string
	Example         string
	Instructions    string
	EditionDate     time.Time
}

func NewExcerciseModifyResponseDTO(excercise models.Excercise) *ExcerciseModifyResponseDTO {
	return &ExcerciseModifyResponseDTO{
		Name:            excercise.Name,
		Description:     excercise.Description,
		CreatorUserID:   utils.GetStringIDFromObjectID(excercise.CreatorUserID),
		Category:        string(excercise.Category),
		MainMuscleGroup: excercise.MainMuscleGroup,
		DifficultLevel:  excercise.DifficultLevel,
		Example:         excercise.Example,
		Instructions:    excercise.Instructions,
		EditionDate:     excercise.EditionDate,
	}
}
