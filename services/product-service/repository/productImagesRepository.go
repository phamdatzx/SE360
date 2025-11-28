package repository

import (
	"context"
	"product-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductImagesRepository interface {
	Create(productImages *model.ProductImages) error
	FindByID(id string) (*model.ProductImages, error)
	FindByProductID(productID string) ([]model.ProductImages, error)
	FindAll() ([]model.ProductImages, error)
	Update(productImages *model.ProductImages) error
	Delete(id string) error
}

type productImagesRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewProductImagesRepository(db *mongo.Database) ProductImagesRepository {
	return &productImagesRepository{
		db:         db,
		collection: db.Collection("product_images"),
	}
}

func (r *productImagesRepository) Create(productImages *model.ProductImages) error {
	productImages.BeforeCreate()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, productImages)
	return err
}

func (r *productImagesRepository) FindByID(id string) (*model.ProductImages, error) {
	var productImages model.ProductImages
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&productImages)
	return &productImages, err
}

func (r *productImagesRepository) FindByProductID(productID string) ([]model.ProductImages, error) {
	var productImages []model.ProductImages
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{"product_id": productID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &productImages); err != nil {
		return nil, err
	}
	return productImages, nil
}

func (r *productImagesRepository) FindAll() ([]model.ProductImages, error) {
	var productImages []model.ProductImages
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &productImages); err != nil {
		return nil, err
	}
	return productImages, nil
}

func (r *productImagesRepository) Update(productImages *model.ProductImages) error {
	productImages.UpdatedAt = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": productImages.ID},
		bson.M{"$set": productImages},
	)
	return err
}

func (r *productImagesRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
