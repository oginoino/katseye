package handlers

import (
	"errors"

	"katseye/internal/domain/services"
	valueobjects "katseye/internal/domain/value_objects"
	"katseye/internal/infrastructure/web/dto"
	"katseye/internal/infrastructure/web/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductHandler struct {
	productService  *services.ProductService
	templateService *services.ProductTemplateService
}

func NewProductHandler(productService *services.ProductService, templateService *services.ProductTemplateService) *ProductHandler {
	return &ProductHandler{
		productService:  productService,
		templateService: templateService,
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
		response.NewInternalServerErrorResponse(c, "Failed to retrieve product", err.Error())
		return
	}

	if product == nil {
		response.NewNotFoundResponse(c, "Product not found", "Product with the given ID does not exist")
		return
	}

	response.NewSuccessResponse(c, "Product retrieved successfully", dto.NewProductResponse(product))
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewBadRequestResponse(c, "Invalid request payload", err.Error())
		return
	}

	product, err := req.ToEntity(primitive.NilObjectID)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid product payload", err.Error())
		return
	}

	if err := h.productService.CreateProduct(c.Request.Context(), product); err != nil {
		switch {
		case errors.Is(err, services.ErrPartnerNotFound):
			response.NewNotFoundResponse(c, "Partner not found", err.Error())
		case errors.Is(err, services.ErrProductTypeNotAccepted):
			response.NewBadRequestResponse(c, "Product type not accepted", err.Error())
		case errors.Is(err, services.ErrPartnerRepositoryUnavailable):
			response.NewInternalServerErrorResponse(c, "Partner data unavailable", err.Error())
		default:
			response.NewInternalServerErrorResponse(c, "Failed to create product", err.Error())
		}
		return
	}

	response.NewCreatedResponse(c, "Product created successfully", dto.NewProductResponse(product))
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid product ID", err.Error())
		return
	}

	var req dto.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewBadRequestResponse(c, "Invalid request payload", err.Error())
		return
	}

	product, err := req.ToEntity(id)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid product payload", err.Error())
		return
	}

	if err := h.productService.UpdateProduct(c.Request.Context(), product); err != nil {
		switch {
		case errors.Is(err, services.ErrPartnerNotFound):
			response.NewNotFoundResponse(c, "Partner not found", err.Error())
		case errors.Is(err, services.ErrProductTypeNotAccepted):
			response.NewBadRequestResponse(c, "Product type not accepted", err.Error())
		case errors.Is(err, services.ErrPartnerRepositoryUnavailable):
			response.NewInternalServerErrorResponse(c, "Partner data unavailable", err.Error())
		default:
			response.NewInternalServerErrorResponse(c, "Failed to update product", err.Error())
		}
		return
	}

	response.NewSuccessResponse(c, "Product updated successfully", dto.NewProductResponse(product))
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid product ID", err.Error())
		return
	}

	existingProduct, err := h.productService.GetProductByID(c.Request.Context(), id)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to retrieve product", err.Error())
		return
	}

	if existingProduct == nil {
		response.NewNotFoundResponse(c, "Product not found", "Product with the given ID does not exist")
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), id); err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to delete product", err.Error())
		return
	}

	response.NewDeleteSuccessResponse(c, "Product", id.Hex())
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	filter := make(map[string]interface{})

	products, err := h.productService.ListProducts(c.Request.Context(), filter)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to retrieve products", err.Error())
		return
	}

	response.NewSuccessResponse(c, "Products retrieved successfully", dto.NewProductResponseList(products))
}

func (h *ProductHandler) ListProductTemplates(c *gin.Context) {
	if h == nil || h.templateService == nil {
		response.NewInternalServerErrorResponse(c, "Product template service unavailable", "handler not configured")
		return
	}

	templates := h.templateService.ListTemplates()
	response.NewSuccessResponse(c, "Product templates retrieved successfully", dto.NewProductTemplateListResponse(templates))
}

func (h *ProductHandler) GetProductTemplate(c *gin.Context) {
	if h == nil || h.templateService == nil {
		response.NewInternalServerErrorResponse(c, "Product template service unavailable", "handler not configured")
		return
	}

	requestedType := c.Param("type")
	productType, err := valueobjects.NewProductType(requestedType)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid product type", err.Error())
		return
	}

	template, ok := h.templateService.GetTemplate(productType)
	if !ok {
		response.NewNotFoundResponse(c, "Product template not found", "template not defined for the provided product type")
		return
	}

	response.NewSuccessResponse(c, "Product template retrieved successfully", dto.NewProductTemplateResponse(template))
}
