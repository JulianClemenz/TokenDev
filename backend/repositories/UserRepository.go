package repositories

import (
	"AppFitness/models"
	"AppFitness/utils"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryInterface interface { //contrato que define metodos para manejar usuarios
	GetUsers() ([]models.User, error)
	GetUsersByID(id string) (models.User, error)
	PostUser(user models.User) (*mongo.InsertOneResult, error)
	PutUser(user models.User) (*models.User, error)
	DeleteUser(id int) error
}

type UserRepository struct { //campo para la conexion a la base de datos
	db DB
}

func NewUserRepository(db DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repository UserRepository) PostUser(user models.User) (*mongo.InsertOneResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("users")
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return result, fmt.Errorf("error al insertar el usuario en UserRepository.PostUser(): %v", err)
	}
	return result, err
}

func (repository UserRepository) GetUsers() ([]models.User, error) { //REVISAR Y AJUSTAR LA DEVOLUCION DE ERRORES
	collection := repository.db.GetClient().Database("AppFitness").Collection("users")
	filter := bson.M{} //filtro vacio para traer todos los documentos

	cursor, err := collection.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	var users []models.User
	for cursor.Next(context.Background()) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, fmt.Errorf("error al decodificar el usuario en UserRepository.GetUsers(): %v", err)
		}
		users = append(users, user)
	}

	return users, err
}

func (repository UserRepository) GetUsersByID(id string) (models.User, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("users")
	objectID := utils.GetObjectIDFromStringID(id)

	filter := bson.M{"_id": objectID}

	result := collection.FindOne(context.TODO(), filter)

	var user models.User
	err := result.Decode(&user)
	if err != nil {
		return models.User{}, fmt.Errorf("error al actualizar el usuario en UserRepository.PutUser(): %v", err)
	}

	return user, nil
}

func (repository UserRepository) GetUserByEmail(email string) (models.User, error) { //para recuperar usuario por email en el login y recuperar ids en service
	collection := repository.db.GetClient().Database("AppFitness").Collection("users")
	filter := bson.M{"email": email}

	result := collection.FindOne(context.TODO(), filter)

	var user models.User
	err := result.Decode(&user)
	if err != nil {
		return models.User{}, fmt.Errorf("error al obtener el usuario en UserRepository.GetUserByEmail(): %v", err)
	}
	return user, nil
}

func (repository UserRepository) PutUser(user models.User) (*mongo.UpdateResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("users")
	filter := bson.M{"_id": user.ID}

	entity := bson.M{"$set": user}
	result, err := collection.UpdateOne(context.TODO(), filter, entity)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repository UserRepository) DeleteUser(id string) (*mongo.DeleteResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("users")
	objectID := utils.GetObjectIDFromStringID(id)
	filter := bson.M{"_id": objectID}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return result, fmt.Errorf("error al eliminar el usuario en UserRepository.DeleteUser(): %v", err)
	}

	return result, err
}
