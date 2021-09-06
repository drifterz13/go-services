package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Find(ctx context.Context) ([]*User, error)
	Create(ctx context.Context, email string) error
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(c *mongo.Collection) UserRepository {
	return &userRepository{c}
}

func (r *userRepository) Find(ctx context.Context) ([]*User, error) {
	cur, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var users []*User
	for cur.Next(ctx) {
		var user User
		if err := cur.Decode(&user); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (r *userRepository) Create(ctx context.Context, email string) error {
	user := &User{
		ID:        primitive.NewObjectID(),
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := r.collection.InsertOne(ctx, &user)

	return err
}
