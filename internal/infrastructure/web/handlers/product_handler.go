package handlers

import (
	"katseye/internal/domain/entities"
	"katseye/internal/domain/services"
	"katseye/internal/infrastructure/web/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid product ID", err.Error())
		return
	}

	product, err := h.productService.GetProductByID(c.Request.Context(), id)
	if err != nil {
		response.NewNotFoundResponse(c, "Product not found", err.Error())
		return
	}

	if product == nil {
		response.NewNotFoundResponse(c, "Product not found", "Product with the given ID does not exist")
		return
	}

	response.NewSuccessResponse(c, "Product retrieved successfully", product)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product entities.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		response.NewBadRequestResponse(c, "Invalid request payload", err.Error())
		return
	}

	err := h.productService.CreateProduct(c.Request.Context(), &product)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to create product", err.Error())
		return
	}

	response.NewCreatedResponse(c, "Product created successfully", product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid product ID", err.Error())
		return
	}

	var product entities.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		response.NewBadRequestResponse(c, "Invalid request payload", err.Error())
		return
	}
	product.ID = id

	err = h.productService.UpdateProduct(c.Request.Context(), &product)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to update product", err.Error())
		return
	}

	response.NewSuccessResponse(c, "Product updated successfully", product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid product ID", err.Error())
		return
	}

	// Optional: Check if product exists before deleting
	existingProduct, err := h.productService.GetProductByID(c.Request.Context(), id)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to retrieve product", err.Error())
		return
	}

	if existingProduct == nil {
		response.NewNotFoundResponse(c, "Product not found", "Product with the given ID does not exist")
		return
	}

	err = h.productService.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to delete product", err.Error())
		return
	}

	response.NewDeleteSuccessResponse(c, "Product", id.Hex())
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	// Extract query parameters for filtering
	filter := make(map[string]interface{})

	// You can add query parameter parsing here
	// Example: if c.Query("category") != "" { filter["product_category"] = c.Query("category") }

	products, err := h.productService.ListProducts(c.Request.Context(), filter)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to retrieve products", err.Error())
		return
	}

	response.NewSuccessResponse(c, "Products retrieved successfully", products)
}
