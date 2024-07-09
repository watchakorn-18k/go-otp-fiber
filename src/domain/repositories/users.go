package repositories

import (
	"context"
	. "go-opt-fiber/src/domain/datasources"
	"go-opt-fiber/src/domain/entities"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type usersRepository struct {
	Context    context.Context
	Collection *mongo.Collection
}

type IUsersRepository interface {
	SaveOrUpdateUser(username, secret string, url string) error
	InsertUser(username, secret string) error
	FindUserByUsername(username string) (*entities.User, error)
}

func NewUsersRepository(db *MongoDB) IUsersRepository {
	return &usersRepository{
		Context:    db.Context,
		Collection: db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("users"),
	}
}

func (r *usersRepository) SaveOrUpdateUser(username, secret string, url string) error {
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"secret": secret, "url": url}}
	opts := options.Update().SetUpsert(true)

	_, err := r.Collection.UpdateOne(r.Context, filter, update, opts)
	return err
}

func (r *usersRepository) InsertUser(username, secret string) error {
	_, err := r.Collection.InsertOne(r.Context, entities.User{
		Username: username,
		Secret:   secret,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *usersRepository) FindUserByUsername(username string) (*entities.User, error) {
	var user entities.User
	filter := bson.M{"username": username}
	if err := r.Collection.FindOne(r.Context, filter).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
