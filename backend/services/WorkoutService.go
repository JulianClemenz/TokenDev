package services

import (
	"AppFitness/dto"
	"AppFitness/models"
	"AppFitness/repositories"
	"fmt"
	"sort"
	"time"
)

type WorkoutInterface interface {
	PostWorkout(*dto.WorkoutRegisterDTO) (*dto.WorkoutResponseDTO, error)
	GetWorkouts(*dto.WorkoutRegisterDTO) ([]*dto.WorkoutResponseDTO, error)
	GetWorkoutByID(id string) (*dto.WorkoutResponseDTO, error)
	DeleteWorkout(id string) error
	GetWorkoutStats(userID string) (*dto.WorkoutStatsDTO, error)
}

type WorkoutService struct {
	WorkoutRepository repositories.WorkoutRepositoryInterface
	RoutineRepository repositories.RoutineRepositoryInterface
	UserRepository    repositories.UserRepositoryInterface
}

func NewWorkoutService(workoutRepository repositories.WorkoutRepositoryInterface) *WorkoutService {
	return &WorkoutService{
		WorkoutRepository: workoutRepository,
	}
}

func (ws WorkoutService) PostWorkout(workoutDTO dto.WorkoutRegisterDTO /*UserID se setea en handler*/) (*dto.WorkoutResponseDTO, error) {
	result, err := ws.RoutineRepository.GetRoutineByID(workoutDTO.RoutineID)
	if err != nil {
		return nil, fmt.Errorf("rutina no encontrada: %w", err)
	}
	if result.ID.IsZero() {
		return nil, fmt.Errorf("rutina no encontrada")
	}

	workoutModel := dto.GetModelWorkoutRegisterDTO(&workoutDTO)
	workoutModel.Date = time.Now()
	workoutModel.RoutineName = result.Name

	insertResult, err := ws.WorkoutRepository.PostWorkout(workoutModel)
	if err != nil {
		return nil, fmt.Errorf("error al crear el workout: %w", err)
	}
	if insertResult.InsertedID == nil {
		return nil, fmt.Errorf("no se pudo crear el workout")
	}

	//recuperar workout creado
	createdWorkout, err := ws.WorkoutRepository.GetWorkoutByID(insertResult.InsertedID.(string))
	if err != nil {
		return nil, fmt.Errorf("error al obtener el workout creado: %w", err)
	}
	if createdWorkout.ID.IsZero() {
		return nil, fmt.Errorf("workout creado no encontrado")
	}

	//convertir a dto y devolver
	workoutResponse := dto.NewWorkoutResponseDTO(createdWorkout)

	return workoutResponse, nil
}

// GetWorkouts obtiene todos los workouts de un usuario específico
func (ws WorkoutService) GetWorkouts(workoutDTO dto.WorkoutRegisterDTO /*UserID se setea en handler*/) ([]*dto.WorkoutResponseDTO, error) {

	//validar existencia de user
	user, err := ws.UserRepository.GetUsersByID(workoutDTO.UserID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario: %w", err)
	}
	if user.ID.IsZero() {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	//obtener workouts del user
	workoutsModel, err := ws.WorkoutRepository.GetWorkoutsByUserID(workoutDTO.UserID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener workouts: %w", err)
	}
	if len(workoutsModel) == 0 {
		return nil, fmt.Errorf("no se encontraron workouts para el usuario")
	}

	var workoutsDTO []*dto.WorkoutResponseDTO
	for _, workout := range workoutsModel {
		workoutDTO := dto.NewWorkoutResponseDTO(workout)
		workoutsDTO = append(workoutsDTO, workoutDTO)
	}

	return workoutsDTO, nil
}

// GetWorkoutByID obtiene un workout por su ID
func (ws WorkoutService) GetWorkoutByID(id string) (*dto.WorkoutResponseDTO, error) {

	//validacion de existencia de user
	user, err := ws.UserRepository.GetUsersByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario: %w", err)
	}
	if user.ID.IsZero() {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	//obtener workout por id
	workoutModel, err := ws.WorkoutRepository.GetWorkoutByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener workout: %w", err)
	}
	if workoutModel.ID.IsZero() {
		return nil, fmt.Errorf("workout no encontrado")
	}
	workoutDTO := dto.NewWorkoutResponseDTO(workoutModel)

	return workoutDTO, nil
}

func (ws WorkoutService) DeleteWorkout(id string) error {
	//validacion de existencia de workout
	workout, err := ws.WorkoutRepository.GetWorkoutByID(id)
	if err != nil {
		return fmt.Errorf("error al obtener workout: %w", err)
	}
	if workout.ID.IsZero() {
		return fmt.Errorf("workout no encontrado")
	}
	//eliminar workout
	result, err := ws.WorkoutRepository.DeleteWorkout(id)
	if err != nil {
		return fmt.Errorf("error al eliminar workout: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no se pudo eliminar el workout")
	}
	return nil
}

func (ws WorkoutService) GetWorkoutStats(userID string) (*dto.WorkoutStatsDTO, error) {

	//validacion de existencia de user
	userModel, err := ws.UserRepository.GetUsersByID(userID)
	if err != nil {
		return nil, fmt.Errorf("Error al recuperar usuario")
	}
	if userModel.ID.IsZero() {
		return nil, fmt.Errorf("No se encontro user")
	}

	//logica
	workoutsUser, err := ws.WorkoutRepository.GetWorkoutsByUserID(userID) //lista de workouts del user
	if err != nil {
		return nil, fmt.Errorf("Error al obtener workouts")
	}
	if len(workoutsUser) == 0 {
		return nil, fmt.Errorf("Error, lista de workouts de user vacia")
	}
	if len(workoutsUser) == 1 {
		return nil, fmt.Errorf("No se pueden calcular estadisticas, solo hay un workout")
	}

	// TotalWorkous
	var status dto.WorkoutStatsDTO
	status.TotalWorkouts = len(workoutsUser)

	//WeeklyFrequency (logica: (cantidad de dias entre el primer y el ult workots - 1) / cantidad de entrenamientos)
	sort.Slice(workoutsUser, func(i, j int) bool { //esta estructura ordena ascendentemente los workouts por fecha de creacion
		return workoutsUser[i].Date.Before(workoutsUser[j].Date)
	})

	var primerWorkout models.Workout
	var ultWorkout models.Workout

	for i, workout := range workoutsUser {
		if i == 0 {
			primerWorkout = workout //obtenemos primer workout
		}
		ultWorkout = workout //obtenemos ultimo workout
	}

	dayDifference := ultWorkout.Date.Sub(primerWorkout.Date).Hours() / 24 //calculamos dias de diferencia
	frequency := (dayDifference - 1) / float64(status.TotalWorkouts)      //calculamos la frecuencia
	status.WeeklyFrequency = frequency

	//MostUsedRoutines (ranking de rutinas mas usadas)
	counts := make(map[string]int, len(workoutsUser)) //agrupa
	for _, w := range workoutsUser {
		counts[w.RoutineName]++ //suma en el map
	}

	//mapear a dto
	out := make([]dto.RoutineUsageDTO, 0, len(counts))
	for name, c := range counts {
		out = append(out, dto.RoutineUsageDTO{
			RoutineName: name,
			Count:       c,
		})
	}
}

func BuildRoutineUsageByName(workouts *[]models.Workout) []dto.RoutineUsageDTO {

	// 1) Crear mapa nombre → cantidad
	counts := make(map[string]int, len(workouts))
	for _, w := range workouts {
		if w.RoutineName == "" {
			continue // opcional: ignorar sin nombre
		}
		counts[w.RoutineName]++
	}

	// 2) Pasar a slice
	out := make([]RoutineUsageDTO, 0, len(counts))
	for name, c := range counts {
		out = append(out, RoutineUsageDTO{
			RoutineName: name,
			Count:       c,
		})
	}

	// 3) Ordenar: más usados primero, luego alfabéticamente
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count == out[j].Count {
			return out[i].RoutineName < out[j].RoutineName
		}
		return out[i].Count > out[j].Count
	})

	return out
}
