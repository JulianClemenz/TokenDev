package repositories

import (
	"AppFitness/models"
	"AppFitness/utils"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExcerciseRepositoryInterface interface {
	PostExcercise(excercise models.Excercise) (*mongo.InsertOneResult, error)
	GetExcercises() ([]models.Excercise, error)
	GetExcerciseByID(id string) (models.Excercise, error)
	PutExcercise(excercise models.Excercise) (*mongo.UpdateResult, error)
	DeleteExcercise(id string) (*mongo.DeleteResult, error)
	ExistByName(name string) (bool, error)
}

type ExcerciseRepository struct { //campo para la conexion a la base de datos
	db DB
}

func NewExcerciseRepository(db DB) *ExcerciseRepository {
	return &ExcerciseRepository{
		db: db,
	}
}

func (repositori ExcerciseRepository) PostExcercise(excercise models.Excercise) (*mongo.InsertOneResult, error) {
	collection := repositori.db.GetClient().Database("AppFitness").Collection("excercises")
	resultado, err := collection.InsertOne(context.TODO(), excercise)
	if err != nil {
		return resultado, fmt.Errorf("error al insertar el ejercicio en ExcerciseRepository.PostExcercise(): %v", err)
	}
	return resultado, err
}

func (repository ExcerciseRepository) GetExcercises() ([]models.Excercise, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("excercises")
	filtro := bson.M{} //filtro vacio para traer todos los documentos

	cursor, err := collection.Find(context.TODO(), filtro)

	defer cursor.Close(context.Background())

	var excercises []models.Excercise
	for cursor.Next(context.Background()) {
		var excercise models.Excercise
		err := cursor.Decode(&excercise)
		if err != nil {
			return nil, fmt.Errorf("error al decodificar el ejercicio en ExcerciseRepository.GetExcercises(): %v", err)
		}
		excercises = append(excercises, excercise)
	}

	return excercises, err
}

func (repository ExcerciseRepository) GetExcerciseByID(id string) (models.Excercise, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("excercises")
	objectID := utils.GetObjectIDFromStringID(id)
	filtro := bson.M{"_id": objectID}

	result := collection.FindOne(context.TODO(), filtro)
	var excercise models.Excercise

	err := result.Decode(&excercise)
	if err != nil {
		return models.Excercise{}, fmt.Errorf("error al obtener el ejercicio en ExcerciseRepository.GetExcerciseByID(): %v", err)
	}
	return excercise, nil
}

func (repository ExcerciseRepository) PutExcercise(excercise models.Excercise) (*mongo.UpdateResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("excercises")
	filtro := bson.M{"_id": excercise.ID}

	entity := bson.M{"$set": bson.M{
		"name":              excercise.Name,
		"description":       excercise.Description,
		"category":          excercise.Category,
		"main_muscle_group": excercise.MainMuscleGroup,
		"example":           excercise.Example,
		"instructions":      excercise.Instructions,
		"edition_date":      excercise.EditionDate,
	}}

	result, err := collection.UpdateOne(context.TODO(), filtro, entity)
	if err != nil {
		return result, fmt.Errorf("error al actualizar el ejercicio en ExcerciseRepository.PutExcercise(): %v", err)
	}
	return result, err
}

func (repository ExcerciseRepository) DeleteExcercise(id string) (*mongo.DeleteResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("excercises")
	objectID := utils.GetObjectIDFromStringID(id)
	filtro := bson.M{"_id": objectID}

	result, err := collection.DeleteOne(context.TODO(), filtro)
	if err != nil {
		return result, fmt.Errorf("error al eliminar el ejercicio en ExcerciseRepository.DeleteExcercise(): %v", err)
	}
	return result, err
}

func (r ExcerciseRepository) ExistByName(name string) (bool, error) {
	collection := r.db.GetClient().Database("AppFitness").Collection("excercises")
	filter := bson.M{"name": name}

	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, fmt.Errorf("error al contar documentos en ExcerciseRepository.ExistByName(): %v", err)
	}

	return count > 0, err
}
