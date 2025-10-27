package services

import (
	"AppFitness/dto"
	"AppFitness/repositories"
	"AppFitness/utils"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoutineInterface interface {
	PostRoutine(routineDTO *dto.RoutineRegisterDTO) (*dto.RoutineResponseDTO, error)
	GetRoutines() ([]*dto.RoutineResponseDTO, error)
	GetRoutineByID(id string) (*dto.RoutineResponseDTO, error)
	PutRoutine(name string, id string) (*dto.RoutineResponseDTO, error)
	AddExcerciseToRoutine(routineID string, exercise *dto.ExcerciseInRoutineDTO) (*dto.RoutineResponseDTO, error)
	RemoveExcerciseFromRoutine(routineID string, excerciseID string) (*dto.RoutineResponseDTO, error)
	UpdateExerciseInRoutine(routineID string, exerciseId string, exerciseMod *dto.ExcerciseInRoutineDTO) (*dto.RoutineResponseDTO, error)
	DeleteRoutine(id string) (bool, error)
}

type RoutineService struct {
	RoutineRepository   repositories.RoutineRepositoryInterface
	ExcerciseRepository repositories.ExcerciseRepositoryInterface
}

func NewRoutineService(routineRepository repositories.RoutineRepositoryInterface, excerciseRspository repositories.ExcerciseRepositoryInterface) *RoutineService {
	return &RoutineService{
		RoutineRepository:   routineRepository,
		ExcerciseRepository: excerciseRspository,
	}
}

func (service *RoutineService) PostRoutine(routineDTO *dto.RoutineRegisterDTO) (*dto.RoutineResponseDTO, error) {
	routineDTO.Name = strings.ToLower(strings.TrimSpace(routineDTO.Name))

	//validacion de campos obligatorios
	if routineDTO.Name == "" {
		return nil, fmt.Errorf("el nombre de la rutina no puede estar vacío")
	}
	if routineDTO.CreatorUserID == "" {
		return nil, fmt.Errorf("el ID del usuario creador no puede estar vacío")
	}

	verificacion, err := service.RoutineRepository.ExistByRutineName(routineDTO.Name)
	if err != nil {
		return nil, fmt.Errorf("no se pudo verificar si existe una rutina con el mismo nombre: %w", err)
	}

	if verificacion {
		return nil, fmt.Errorf("dicho nombre de rutina ya existe")
	}

	//LOGICA
	model := dto.GetModelRoutineRegisterDTO(routineDTO)
	result, err := service.RoutineRepository.PostRoutine(model) //insertamos la rutina en la base de datos
	if err != nil {
		return nil, fmt.Errorf("error al crear la rutina en RoutineService.PostRoutine(): %v", err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("no se pudo obtener el ObjectID insertado")
	}
	idStr := utils.GetStringIDFromObjectID(oid)

	routineModel, err := service.RoutineRepository.GetRoutineByID(idStr) //obtenemos la rutina de tipo response creada para devolver

	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina creada en RoutineService.PostRoutine(): %v", err)
	}

	return dto.NewRoutineResponseDTO(routineModel), nil
}

func (service *RoutineService) GetRoutines() ([]*dto.RoutineResponseDTO, error) {
	routinesDB, err := service.RoutineRepository.GetRoutines()
	if err != nil {
		return nil, fmt.Errorf("error al obtener rutinas %v:", err)
	}
	var routines []*dto.RoutineResponseDTO
	for _, routineDB := range routinesDB {
		routine := dto.NewRoutineResponseDTO(routineDB)
		routines = append(routines, routine)
	}

	return routines, nil
}

func (service *RoutineService) GetRoutineByID(id string) (*dto.RoutineResponseDTO, error) {
	routineDB, err := service.RoutineRepository.GetRoutineByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina por ID en RoutineService.GetRoutineByID(): %v", err)
	}
	return dto.NewRoutineResponseDTO(routineDB), nil
}

// PUT SOLO DE NAME
func (service *RoutineService) PutRoutine(name string, id string) (*dto.RoutineResponseDTO, error) {

	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("el nombre de la rutina no puede estar vacío")
	}
	routineDB, err := service.RoutineRepository.GetRoutineByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina a modificar en RoutineService.PutRoutine(): %v", err)
	}
	if routineDB.Name == name {
		return nil, fmt.Errorf("el nuevo nombre de la rutina no puede ser igual al anterior")
	}

	routineDB.Name = strings.ToLower(strings.TrimSpace(name))
	routineDB.EditionDate = time.Now()

	result, err := service.RoutineRepository.PutRoutine(routineDB)
	if err != nil {
		return nil, fmt.Errorf("error al modificar la rutina en RoutineService.PutRoutine(): %v", err)
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no se modificó ninguna rutina")
	}
	updatedRoutineDB, err := service.RoutineRepository.GetRoutineByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina modificada en RoutineService.PutRoutine(): %v", err)
	}
	return dto.NewRoutineResponseDTO(updatedRoutineDB), nil
}

func (service *RoutineService) AddExcerciseToRoutine(routineID string, exercise *dto.ExcerciseInRoutineDTO) (*dto.RoutineResponseDTO, error) {

	//busqueda de ej
	exerciseDB, err := service.ExcerciseRepository.GetExcerciseByID(exercise.ExcerciseID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el ejercicio a agregar: %w", err)
	}
	if exerciseDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ningún ejercicio con ese ID")
	}

	//busqueda de rutina
	routineDB, err := service.RoutineRepository.GetRoutineByID(routineID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina a modificar: %w", err)
	}
	if routineDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ninguna rutina con ese ID")
	}
	id := utils.GetObjectIDFromStringID(routineID) //convertimos para pasar a el repository

	//convertimos dto a model
	exerciseModel := dto.GetModelExerciseInRoutineDTO(exercise)
	exerciseModel.CreationDate = time.Now()

	//agregar ejercicio a la rutina
	result, err := service.RoutineRepository.AddExerciseRutine(exerciseModel, id)
	if err != nil {
		return nil, fmt.Errorf("error al agregar el ejercicio a la rutina: %w", err)
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no se agregó ningún ejercicio a la rutina")
	}

	//modificar la fecha de edición de la rutina
	routineDB.EditionDate = time.Now()
	_, err = service.RoutineRepository.PutRoutine(routineDB)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar la fecha de edición de la rutina en RoutineService.AddExcerciseToRoutine(): %v", err)
	}

	//buscamos rutina para devolver
	updatedRoutineDB, err := service.RoutineRepository.GetRoutineByID(routineID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina modificada en RoutineService.AddExcerciseToRoutine(): %v", err)
	}
	if updatedRoutineDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ninguna rutina con ese ID, error al agregar ejercicio")
	}
	return dto.NewRoutineResponseDTO(updatedRoutineDB), nil
}

func (service *RoutineService) RemoveExcerciseFromRoutine(routineID string, excerciseID string) (*dto.RoutineResponseDTO, error) {
	//validaciones
	routineDB, err := service.RoutineRepository.GetRoutineByID(routineID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina a modificar: %w", err)
	}
	if routineDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ninguna rutina con ese ID")
	}
	routineObjectID := utils.GetObjectIDFromStringID(routineID)

	exerciseDB, err := service.ExcerciseRepository.GetExcerciseByID(excerciseID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el ejercicio a eliminar: %w", err)
	}
	if exerciseDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ningún ejercicio con ese ID")
	}
	exerciseObjectID := utils.GetObjectIDFromStringID(excerciseID)

	//lógica de eliminación
	result, err := service.RoutineRepository.DeleteExerciseToRutine(routineObjectID, exerciseObjectID)
	if err != nil {
		return nil, fmt.Errorf("error al eliminar el ejercicio de la rutina: %w", err)
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no se eliminó ningún ejercicio de la rutina")
	}

	//modificar la fecha de edición de la rutina
	routineDB.EditionDate = time.Now()
	_, err = service.RoutineRepository.PutRoutine(routineDB)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar la fecha de edición de la rutina en RoutineService.RemoveExcerciseFromRoutine(): %v", err)
	}

	//buscamos rutina para devolver
	updatedRoutineDB, err := service.RoutineRepository.GetRoutineByID(routineID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina modificada en RoutineService.RemoveExcerciseFromRoutine(): %v", err)
	}
	if updatedRoutineDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ninguna rutina con ese ID, error al eliminar ejercicio")
	}
	return dto.NewRoutineResponseDTO(updatedRoutineDB), nil
}

func (service *RoutineService) UpdateExerciseInRoutine(routineID string, exerciseId string, exerciseMod *dto.ExcerciseInRoutineDTO) (*dto.RoutineResponseDTO, error) {

	//validaciones
	routineDB, err := service.RoutineRepository.GetRoutineByID(routineID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina a modificar: %w", err)
	}
	if routineDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ninguna rutina con ese ID")
	}
	routineObjectID := utils.GetObjectIDFromStringID(routineID)

	exerciseDB, err := service.ExcerciseRepository.GetExcerciseByID(exerciseId)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el ejercicio a modificar: %w", err)
	}
	if exerciseDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ningún ejercicio con ese ID")
	}
	exerciseObjectID := utils.GetObjectIDFromStringID(exerciseId)

	//lógica de modificación
	exerciseModel := dto.GetModelExerciseInRoutineDTO(exerciseMod)
	result, err := service.RoutineRepository.UpdateExerciseInRoutine(routineObjectID, exerciseObjectID, exerciseModel)
	if err != nil {
		return nil, fmt.Errorf("error al modificar el ejercicio de la rutina: %w", err)
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no se modificó ningún ejercicio de la rutina")
	}

	//modificar la fecha de edición de la rutina
	routineDB.EditionDate = time.Now()
	_, err = service.RoutineRepository.PutRoutine(routineDB)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar la fecha de edición de la rutina en RoutineService.UpdateExerciseInRoutine(): %v", err)
	}

	//buscamos rutina para devolver
	updatedRoutineDB, err := service.RoutineRepository.GetRoutineByID(routineID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la rutina modificada en RoutineService.UpdateExerciseInRoutine(): %v", err)
	}
	if updatedRoutineDB.ID.IsZero() {
		return nil, fmt.Errorf("no existe ninguna rutina con ese ID, error al modificar ejercicio")
	}
	return dto.NewRoutineResponseDTO(updatedRoutineDB), nil
}

func (service *RoutineService) DeleteRoutine(id string) (bool, error) {
	//validacion de existencia
	exist, err := service.RoutineRepository.GetRoutineByID(id)
	if err != nil {
		return false, fmt.Errorf("error al verificar la existencia de la rutina en RoutineService.DeleteRoutine(): %v", err)
	}
	if exist.ID.IsZero() {
		return false, fmt.Errorf("no existe ninguna rutina con ese ID")
	}

	result, err := service.RoutineRepository.DeleteRoutine(id)
	if err != nil {
		return false, fmt.Errorf("error al eliminar la rutina en RoutineService.DeleteRoutine(): %v", err)
	}
	if result.DeletedCount == 0 {
		return false, fmt.Errorf("no se eliminó ninguna rutina")
	}
	return false, nil
}
