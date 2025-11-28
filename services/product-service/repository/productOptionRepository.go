package repository

import (
	"context"
	"product-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductOptionRepository interface {
	Create(productOption *model.ProductOption) error
	FindByID(id string) (*model.ProductOption, error)
	FindByProductID(productID string) ([]model.ProductOption, error)
	FindAll() ([]model.ProductOption, error)
	Update(productOption *model.ProductOption) error
	Delete(id string) error
}

type productOptionRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewProductOptionRepository(db *mongo.Database) ProductOptionRepository {
	return &productOptionRepository{
		db:         db,
		collection: db.Collection("product_options"),
	}
}

func (r *productOptionRepository) Create(productOption *model.ProductOption) error {
	productOption.BeforeCreate()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, productOption)
	return err
}

func (r *productOptionRepository) FindByID(id string) (*model.ProductOption, error) {
	var productOption model.ProductOption
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&productOption)
	return &productOption, err
}

func (r *productOptionRepository) FindByProductID(productID string) ([]model.ProductOption, error) {
	var productOptions []model.ProductOption
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{"product_id": productID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &productOptions); err != nil {
		return nil, err
	}
	return productOptions, nil
}

func (r *productOptionRepository) FindAll() ([]model.ProductOption, error) {
	var productOptions []model.ProductOption
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &productOptions); err != nil {
		return nil, err
	}
	return productOptions, nil
}

func (r *productOptionRepository) Update(productOption *model.ProductOption) error {
	productOption.UpdatedAt = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": productOption.ID},
		bson.M{"$set": productOption},
	)
	return err
}

func (r *productOptionRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
