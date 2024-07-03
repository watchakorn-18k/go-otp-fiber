package domain

import (
	"context"
	"go-opt-fiber/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveOrUpdateUser(username, secret string, url string, qrCode []byte, userCollection *mongo.Collection) error {
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"secret": secret, "url": url, "qrcode": qrCode}}
	opts := options.Update().SetUpsert(true)

	_, err := userCollection.UpdateOne(context.Background(), filter, update, opts)
	return err
}

func InsertUser(username, secret string, userCollection *mongo.Collection) error {
	_, err := userCollection.InsertOne(context.Background(), entities.User{
		Username: username,
		Secret:   secret,
	})
	if err != nil {
		return err
	}
	return nil
}

func FindUserByUsername(username string, userCollection *mongo.Collection) (*entities.User, error) {
	var user entities.User
	err := userCollection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
