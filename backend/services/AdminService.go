package services

import (
	"AppFitness/dto"
	"AppFitness/repositories"
	"AppFitness/utils"
	"fmt"
	"sort"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminInterface interface {
	GetGlobalStats() ([]*dto.TopUsedExcerciseDTO, error)
	GetLogs() ([]*dto.UserResponseDTO, int, error) //lista de users, cantidad
}

type AdminService struct {
	UserRepository      repositories.UserRepository
	ExcerciseRepository repositories.ExcerciseRepository
	Routinerepository   repositories.RoutineRepository
}

func NewAdminService(UserRepository repositories.UserRepository, ExcerciseRepository repositories.ExcerciseRepository) *AdminService {
	return &AdminService{
		UserRepository:      UserRepository,
		ExcerciseRepository: ExcerciseRepository,
	}
}

func (a *AdminService) GetGlobalStats() ([]*dto.TopUsedExcerciseDTO, error) {
	routinesDB, err := a.Routinerepository.GetRoutines()
	if err != nil {
		return nil, fmt.Errorf("error al recuperar rutinas")
	}
	if len(routinesDB) == 0 {
		return []*dto.TopUsedExcerciseDTO{}, fmt.Errorf("vacio")
	}
	//recuperamos todos los exc
	var excerciseList []primitive.ObjectID
	for _, r := range routinesDB {
		for _, e := range r.ExcerciseList {
			excerciseList = append(excerciseList, e.ExcerciseID)
		}
	}
	//agrupamos por id
	counts := make(map[primitive.ObjectID]int, len(excerciseList))
	for _, e := range excerciseList {
		counts[e]++
	}
	//mapear a dto TopUsedExcerciseDTO id, name, count

	topList := make([]*dto.TopUsedExcerciseDTO, 0, len(counts))
	for te, i := range counts {
		excerciseName, _ := a.ExcerciseRepository.GetExcerciseByID(utils.GetStringIDFromObjectID(te))
		topList = append(topList, &dto.TopUsedExcerciseDTO{
			ExcerciseID:   utils.GetStringIDFromObjectID(te),
			ExcerciseName: excerciseName.Name,
			Count:         i,
		})

	}
	//ordenamos primero por mas usados, segundo alfabeticamente
	sort.Slice(topList, func(i, j int) bool {
		if topList[i].Count == topList[j].Count {
			return topList[i].ExcerciseName < topList[j].ExcerciseName
		}
		return topList[i].Count > topList[j].Count
	})
	return topList, nil
}

func (a *AdminService) GetLogs() ([]*dto.UserResponseDTO, int, error) {

	//buscar users
	usersDB, err := a.UserRepository.GetUsers()
	if err != nil {
		return nil, 0, fmt.Errorf("error al recuperar users")
	}
	if len(usersDB) == 0 {
		return []*dto.UserResponseDTO{}, 0, fmt.Errorf("vacio")
	}

	var usersResponse []*dto.UserResponseDTO
	for _, u := range usersDB {
		usersResponse = append(usersResponse, dto.NewUserResponseDTO(u))
	}

	return usersResponse, len(usersResponse), nil

}
