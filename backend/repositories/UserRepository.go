package repositories

import (
	"AppFitness/models"
	"AppFitness/utils"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryInterface interface { //contrato que define metodos para manejar usuarios
	GetUsers() ([]models.User, error)
	GetUSersByID(id string) (models.User, error)
	PostUser(user models.User) (*models.User, error)
	PutUser(id int, user models.User) (*models.User, error)
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
	return result, err
}

func (repository UserRepository) GetUsers() ([]models.User, error) { //REVISAR Y AJUSTAR LA DEVOLUCION DE ERRORES
	collection := repository.db.GetClient().Database("AppFitness").Collection("users")
	filtro := bson.M{} //filtro vacio para traer todos los documentos

	cursor, err := collection.Find(context.TODO(), filtro)
	defer cursor.Close(context.TODO())

	var users []models.User
	for cursor.Next(context.Background()) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, err
}

func (repository UserRepository) GetUSersByID(id string) (models.User, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("users")
	objectID := utils.GetObjectIDFromStringID(id)

	filtro := bson.M{"_id": objectID}

	result := collection.FindOne(context.TODO(), filtro)

	var user models.User
	err := result.Decode(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil

	/*cursor, err := collection.Find(context.TODO(), filtro)
	defer cursor.Close(context.Background())

	//itera a travez de los resultados
	var user models.User
	for cursor.Next(context.Background()) {
		err := cursor.Decode(&user)
		if err != nil {
			return models.User{}, fmt.Errorf("usuario con id %s no encontrado", id)
		}
	}

	return user, err*/
}

func (repository UserRepository) PutUser(user models.User) (*mongo.UpdateResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("users")
	filtro := bson.M{"_id": user.ID}

	entidad := bson.M{"$set": user}
	resultado, err := collection.UpdateOne(context.TODO(), filtro, entidad)

	return resultado, err
}

func (repository UserRepository) DeleteUser(id string) (*mongo.DeleteResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("users")
	objectID := utils.GetObjectIDFromStringID(id)
	filtro := bson.M{"_id": objectID}

	resultado, err := collection.DeleteOne(context.TODO(), filtro)

	return resultado, err
}
