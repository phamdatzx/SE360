package repository

import (
	"context"
	"product-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductReportRepository interface {
	Create(productReport *model.ProductReport) error
	FindByID(id string) (*model.ProductReport, error)
	FindByProductID(productID string) ([]model.ProductReport, error)
	FindByUserID(userID string) ([]model.ProductReport, error)
	FindAll() ([]model.ProductReport, error)
	Update(productReport *model.ProductReport) error
	Delete(id string) error
}

type productReportRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewProductReportRepository(db *mongo.Database) ProductReportRepository {
	return &productReportRepository{
		db:         db,
		collection: db.Collection("product_reports"),
	}
}

func (r *productReportRepository) Create(productReport *model.ProductReport) error {
	productReport.BeforeCreate()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, productReport)
	return err
}

func (r *productReportRepository) FindByID(id string) (*model.ProductReport, error) {
	var productReport model.ProductReport
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&productReport)
	return &productReport, err
}

func (r *productReportRepository) FindByProductID(productID string) ([]model.ProductReport, error) {
	var productReports []model.ProductReport
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{"product_id": productID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &productReports); err != nil {
		return nil, err
	}
	return productReports, nil
}

func (r *productReportRepository) FindByUserID(userID string) ([]model.ProductReport, error) {
	var productReports []model.ProductReport
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &productReports); err != nil {
		return nil, err
	}
	return productReports, nil
}

func (r *productReportRepository) FindAll() ([]model.ProductReport, error) {
	var productReports []model.ProductReport
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &productReports); err != nil {
		return nil, err
	}
	return productReports, nil
}

func (r *productReportRepository) Update(productReport *model.ProductReport) error {
	productReport.UpdatedAt = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": productReport.ID},
		bson.M{"$set": productReport},
	)
	return err
}

func (r *productReportRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
