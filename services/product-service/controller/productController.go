package controller

import (
	"net/http"
	"product-service/error"
	"product-service/model"
	"product-service/service"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service service.ProductService
}

func NewProductController(service service.ProductService) *ProductController {
	return &ProductController{service: service}
}

func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var product model.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Invalid request body", err))
		return
	}

	if err := c.service.CreateProduct(&product); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to create product", err))
		return
	}

	ctx.JSON(http.StatusCreated, product)
}

func (c *ProductController) GetProductByID(ctx *gin.Context) {
	id := ctx.Param("id")

	product, err := c.service.GetProductByID(id)
	if err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusNotFound, "Product not found", err))
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (c *ProductController) GetAllProducts(ctx *gin.Context) {
	products, err := c.service.GetAllProducts()
	if err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to fetch products", err))
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	var product model.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusBadRequest, "Invalid request body", err))
		return
	}
	product.ID = id

	if err := c.service.UpdateProduct(&product); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to update product", err))
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeleteProduct(id); err != nil {
		ctx.Error(error.NewAppErrorWithErr(http.StatusInternalServerError, "Failed to delete product", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
