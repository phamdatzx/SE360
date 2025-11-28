package repository

import (
	"context"
	"product-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SellerCategoryRepository interface {
	Create(sellerCategory *model.SellerCategory) error
	FindByID(id string) (*model.SellerCategory, error)
	FindAll() ([]model.SellerCategory, error)
	Update(sellerCategory *model.SellerCategory) error
	Delete(id string) error
}

type sellerCategoryRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewSellerCategoryRepository(db *mongo.Database) SellerCategoryRepository {
	return &sellerCategoryRepository{
		db:         db,
		collection: db.Collection("seller_categories"),
	}
}

func (r *sellerCategoryRepository) Create(sellerCategory *model.SellerCategory) error {
	sellerCategory.BeforeCreate()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, sellerCategory)
	return err
}

func (r *sellerCategoryRepository) FindByID(id string) (*model.SellerCategory, error) {
	var sellerCategory model.SellerCategory
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&sellerCategory)
	return &sellerCategory, err
}

func (r *sellerCategoryRepository) FindAll() ([]model.SellerCategory, error) {
	var sellerCategories []model.SellerCategory
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &sellerCategories); err != nil {
		return nil, err
	}
	return sellerCategories, nil
}

func (r *sellerCategoryRepository) Update(sellerCategory *model.SellerCategory) error {
	sellerCategory.UpdatedAt = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": sellerCategory.ID},
		bson.M{"$set": sellerCategory},
	)
	return err
}

func (r *sellerCategoryRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
