// Package users implement all logic.
package users

import (
	"context"
	"fmt"
	"time"

	"github.com/kubuskotak/ymir-test/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (i *impl) GetAll(ctx context.Context, request entity.RequestGetUsers) (result entity.ResponseGetUsers, err error) {
	coll := i.adapter.PersistUsers.Collection("users")

	skip := (request.Page - 1) * request.Limit

	// Query options with skip and limit
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(request.Limit))

	// pagination
	result.Limit = request.Limit
	result.Page = request.Page
	var cursor *mongo.Cursor
	cursor, err = coll.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return result, err
	}
	defer func(c context.Context) {
		err = cursor.Close(c)
	}(ctx)

	// Iterate through the cursor to get each document.
	var documents = make([]entity.User, 0)
	for cursor.Next(ctx) {
		var document entity.User
		err := cursor.Decode(&document)
		if err != nil {
			return result, err
		}
		documents = append(documents, document)
	}

	result.Users = documents
	return result, nil
}

func (i *impl) Create(ctx context.Context, user entity.User) (entity.User, error) {
	coll := i.adapter.PersistUsers.Collection("users")

	user.CreatedAt = time.Now()

	result, err := coll.InsertOne(ctx, user)
	if err != nil {
		return entity.User{}, err
	}

	// Retrieve the created document using the _id from the InsertOneResult
	var createdUser entity.User
	err = coll.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&createdUser)
	if err != nil {
		return entity.User{}, err
	}

	return createdUser, nil
}

func (i *impl) GetById(ctx context.Context, userId string) (entity.User, error) {
	coll := i.adapter.PersistUsers.Collection("users")
	var createdUser entity.User

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return entity.User{}, err
	}

	filter := bson.D{{Key: "_id", Value: id}}

	err = coll.FindOne(ctx, filter).Decode(&createdUser)
	if err != nil {
		return entity.User{}, err
	}
	return createdUser, nil
}

func (i *impl) UpdateById(ctx context.Context, user entity.User) (entity.User, error) {
	coll := i.adapter.PersistUsers.Collection("users")

	id, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return entity.User{}, err
	}

	filter := bson.D{{Key: "_id", Value: id}}

	// Attempt to find the document
	var result bson.M
	err = coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return entity.User{}, err
	}

	// The updates
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: user.Name},
			{Key: "email", Value: user.Email},
			{Key: "age", Value: user.Age},
			// Add more fields here if needed
		}},
	}

	_, err = coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return entity.User{}, err
	}

	// Query the updated user data
	err = coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil

}

func (i *impl) DeleteById(ctx context.Context, userId string) error {
	coll := i.adapter.PersistUsers.Collection("users")

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: id}}

	// Attempt to find the document
	var result bson.M
	err = coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return fmt.Errorf("no document with id %v was found", userId)
	}

	// If document is found, attempt to delete
	_, err = coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
