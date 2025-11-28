package repository

import (
	"context"
	"product-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RatingImageRepository interface {
	Create(ratingImage *model.RatingImage) error
	FindByID(id string) (*model.RatingImage, error)
	FindByRatingID(ratingID string) ([]model.RatingImage, error)
	FindAll() ([]model.RatingImage, error)
	Update(ratingImage *model.RatingImage) error
	Delete(id string) error
}

type ratingImageRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewRatingImageRepository(db *mongo.Database) RatingImageRepository {
	return &ratingImageRepository{
		db:         db,
		collection: db.Collection("rating_images"),
	}
}

func (r *ratingImageRepository) Create(ratingImage *model.RatingImage) error {
	ratingImage.BeforeCreate()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, ratingImage)
	return err
}

func (r *ratingImageRepository) FindByID(id string) (*model.RatingImage, error) {
	var ratingImage model.RatingImage
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&ratingImage)
	return &ratingImage, err
}

func (r *ratingImageRepository) FindByRatingID(ratingID string) ([]model.RatingImage, error) {
	var ratingImages []model.RatingImage
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{"rating_id": ratingID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &ratingImages); err != nil {
		return nil, err
	}
	return ratingImages, nil
}

func (r *ratingImageRepository) FindAll() ([]model.RatingImage, error) {
	var ratingImages []model.RatingImage
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &ratingImages); err != nil {
		return nil, err
	}
	return ratingImages, nil
}

func (r *ratingImageRepository) Update(ratingImage *model.RatingImage) error {
	ratingImage.UpdatedAt = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": ratingImage.ID},
		bson.M{"$set": ratingImage},
	)
	return err
}

func (r *ratingImageRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
