package repository

import (
	"context"
	"product-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RatingRepository interface {
	Create(rating *model.Rating) error
	FindByID(id string) (*model.Rating, error)
	FindByUserID(userID string) ([]model.Rating, error)
	FindByParentID(parentID string) ([]model.Rating, error)
	FindAll() ([]model.Rating, error)
	Update(rating *model.Rating) error
	Delete(id string) error
}

type ratingRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewRatingRepository(db *mongo.Database) RatingRepository {
	return &ratingRepository{
		db:         db,
		collection: db.Collection("ratings"),
	}
}

func (r *ratingRepository) Create(rating *model.Rating) error {
	rating.BeforeCreate()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, rating)
	return err
}

func (r *ratingRepository) FindByID(id string) (*model.Rating, error) {
	var rating model.Rating
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&rating)
	return &rating, err
}

func (r *ratingRepository) FindByUserID(userID string) ([]model.Rating, error) {
	var ratings []model.Rating
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &ratings); err != nil {
		return nil, err
	}
	return ratings, nil
}

func (r *ratingRepository) FindByParentID(parentID string) ([]model.Rating, error) {
	var ratings []model.Rating
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{"parent_id": parentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &ratings); err != nil {
		return nil, err
	}
	return ratings, nil
}

func (r *ratingRepository) FindAll() ([]model.Rating, error) {
	var ratings []model.Rating
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &ratings); err != nil {
		return nil, err
	}
	return ratings, nil
}

func (r *ratingRepository) Update(rating *model.Rating) error {
	rating.UpdatedAt = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": rating.ID},
		bson.M{"$set": rating},
	)
	return err
}

func (r *ratingRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
