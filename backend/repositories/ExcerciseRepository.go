package repositories

import (
	"AppFitness/models"
	"context"
)

func ExcerciseRepositoryInterface interface {
	PostExcercise(excercise models.Excercise) (*mongo.InsertOneResult, error)
	GetExcercises() ([]models.Excercise, error)
	GetExcerciseByID(id string) (models.Excercise, error)
	PutExcercise(excercise models.Excercise) (*mongo.UpdateResult, error)
	DeleteExcercise(id string) (*mongo.DeleteResult, error)
}

type ExcerciseRepository struct { //campo para la conexion a la base de datos
	db DB
}

func NewExcerciseRepository(db DB) *ExcerciseRepository {
	return &ExcerciseRepository{
		db: db,
	}
}

func (repositori ExcerciseRepository) PostExcercise(excercise models.Excercise)(*mongo.InsertOneResult, error){
	collection := repositori.db.GetClient().Database("AppFitness").Collection("excercises")
	resultado, err := collection.InsertOne(context.TODO(), excercise)
	return resultado, err
}

func (repository ExcerciseRepository) GetExcercises()([]models.Excercise, error){
	collection := repository.db.GetClient().Database("AppFitness").Collection("excercises")
	filtro := bson.M{} //filtro vacio para traer todos los documentos

	cursor, err := collection.Find(context.TODO(), filtro)

	defer cursor.Close(context.Background())

	var excercises []models.Excercise
	for cursor.Next(context.Background()) {
		var excercise models.Excercise
		err := cursor.Decode(&excercise)
		if err != nil {
			return nil, err
		}
		excercises = append(excercises, excercise)
	}

	return excercises, err
}

func (repository ExcerciseRepository) GetExcerciseByID(id string)(models.Excercise, error){
	collection := repository.db.GetClient().Database("AppFitness").Collection("excercises")
	objectID := utils.GetObjectIDFromStringID(id)
	filtro := bson.M{"_id": objectID}

	cursor, err := collection.FindOne(context.TODO(), filtro)
}