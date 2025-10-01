package repositories

import (
	"AppFitness/models"
	"AppFitness/utils"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoutineRepositoryInterface interface {
	PostRoutine(models.Routine) (*mongo.InsertOneResult, error)
	GetRoutines() ([]models.Routine, error)
	GetRoutineByID(id string) (models.Routine, error)
	PutRoutine(routine models.Routine) (*mongo.UpdateResult, error)
	DeleteRoutine(id string) (*mongo.DeleteResult, error)
}

type RoutineRepository struct {
	db DB
}

func NewRoutineRepository(db DB) *RoutineRepository {
	return &RoutineRepository{
		db: db,
	}
}

func (repository RoutineRepository) PostRoutine(routine models.Routine) (*mongo.InsertOneResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("routines")
	result, err := collection.InsertOne(nil, routine)
	if err != nil {
		return result, fmt.Errorf("error al insertar la rutina en RoutineRepository.PostRoutine(): %v", err)
	}
	return result, nil
}

func (repository RoutineRepository) GetRoutines() ([]models.Routine, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("routines")
	filter := bson.M{}

	cursor, err := collection.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	var routines []models.Routine
	for cursor.Next(context.Background()) {
		var routine models.Routine
		err := cursor.Decode(&routine)
		if err != nil {
			return nil, fmt.Errorf("error al decodificar la rutina en RoutineRepository.GetRoutines(): %v", err)
		}
		routines = append(routines, routine)
	}
	return routines, err
}

func (repository RoutineRepository) GetRoutineByID(id string) (models.Routine, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("routines")
	objectID := utils.GetObjectIDFromStringID(id)
	filter := bson.M{"_id": objectID}

	result := collection.FindOne(context.TODO(), filter)

	var routine models.Routine
	err := result.Decode(&routine)
	if err != nil {
		return models.Routine{}, fmt.Errorf("error al obtener la rutina en RoutineRepository.GetRoutineByID(): %v", err)
	}
	return routine, nil
}

func (repository RoutineRepository) PutRoutine(routine models.Routine) (*mongo.UpdateResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("routines")
	filter := bson.M{"_id": routine.ID}
	entity := bson.M{"$set": routine}
	result, err := collection.UpdateOne(context.TODO(), filter, entity)
	if err != nil {
		return result, fmt.Errorf("error al actualizar la rutina en RoutineRepository.PutRoutine(): %v", err)
	}
	return result, nil
}

func (repository RoutineRepository) DeleteRoutine(id string) (*mongo.DeleteResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("routines")
	objectID := utils.GetObjectIDFromStringID(id)
	filter := bson.M{"_id": objectID}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return result, fmt.Errorf("error al eliminar la rutina en RoutineRepository.DeleteRoutine(): %v", err)
	}
	return result, nil
}
