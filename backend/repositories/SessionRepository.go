package repositories

import (
	"AppFitness/models"
	"AppFitness/utils"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionRepositoryInterface interface {
	PostSession(models.Session) (*mongo.InsertOneResult, error)
	GetSessions() ([]models.Session, error)
	GetSessionByID(id string) (models.Session, error)
	PutSession(session models.Session) (*mongo.UpdateResult, error)
	DeleteSession(id string) (*mongo.DeleteResult, error)
}

type SessionRepository struct {
	db DB
}

func NewSessionRepository(db DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (repository SessionRepository) PostSession(session models.Session) (*mongo.InsertOneResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("sessions")
	result, err := collection.InsertOne(nil, session)
	if err != nil {
		return result, fmt.Errorf("error al insertar la session en SessionRepository.PostSession(): %v", err)
	}
	return result, nil
}

func (repository SessionRepository) GetSessions() ([]models.Session, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("sessions")
	filter := bson.M{}

	cursor, err := collection.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())
	var sessions []models.Session
	for cursor.Next(context.Background()) {
		var session models.Session
		err := cursor.Decode(&session)
		if err != nil {
			return nil, fmt.Errorf("error al decodificar la session en SessionRepository.GetSessions(): %v", err)
		}
		sessions = append(sessions, session)
	}
	return sessions, err
}

func (repository SessionRepository) GetSessionByID(id string) (models.Session, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("sessions")
	objectID := utils.GetObjectIDFromStringID(id)
	filter := bson.M{"_id": objectID}

	result := collection.FindOne(context.TODO(), filter)
	var session models.Session
	err := result.Decode(&session)
	if err != nil {
		return models.Session{}, fmt.Errorf("error al obtener la session en SessionRepository.GetSessionByID(): %v", err)
	}
	return session, nil
}

func (repository SessionRepository) PutSession(session models.Session) (*mongo.UpdateResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("sessions")
	objectID := utils.GetObjectIDFromStringID(session.ID.Hex())
	filtro := bson.M{"_id": objectID}
	entity := bson.M{"$set": session}
	result, err := collection.UpdateOne(context.TODO(), filtro, entity)
	if err != nil {
		return result, fmt.Errorf("error al obtener la session SessionRepository.PutSession(): %v", err)
	}
	return result, nil
}

func (repository SessionRepository) DeleteSession(id string) (*mongo.DeleteResult, error) {
	collection := repository.db.GetClient().Database("AppFitness").Collection("sessions")
	objectID := utils.GetObjectIDFromStringID(id)
	filter := bson.M{"_id": objectID}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return result, fmt.Errorf("error al eliminar la session en SessionRepository.DeleteSession(): %v", err)
	}
	return result, nil
}