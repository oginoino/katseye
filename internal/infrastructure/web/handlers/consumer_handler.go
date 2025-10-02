package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"katseye/internal/domain/entities"
	"katseye/internal/domain/services"
	"katseye/internal/infrastructure/web/dto"
	"katseye/internal/infrastructure/web/response"
)

type ConsumerHandler struct {
	consumerService *services.ConsumerService
}

func NewConsumerHandler(consumerService *services.ConsumerService) *ConsumerHandler {
	return &ConsumerHandler{consumerService: consumerService}
}

func (h *ConsumerHandler) GetConsumer(c *gin.Context) {
	if h == nil || h.consumerService == nil {
		response.NewInternalServerErrorResponse(c, "Consumer service unavailable", "consumer service not configured")
		return
	}

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid consumer ID", err.Error())
		return
	}

	consumer, err := h.consumerService.GetConsumerByID(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrConsumerRepositoryUnavailable):
			response.NewInternalServerErrorResponse(c, "Consumer data unavailable", err.Error())
		default:
			response.NewInternalServerErrorResponse(c, "Failed to retrieve consumer", err.Error())
		}
		return
	}

	if consumer == nil {
		response.NewNotFoundResponse(c, "Consumer not found", "Consumer with the given ID does not exist")
		return
	}

	response.NewSuccessResponse(c, "Consumer retrieved successfully", dto.NewConsumerResponse(consumer))
}

func (h *ConsumerHandler) CreateConsumer(c *gin.Context) {
	if h == nil || h.consumerService == nil {
		response.NewInternalServerErrorResponse(c, "Consumer service unavailable", "consumer service not configured")
		return
	}

	var req dto.ConsumerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewBadRequestResponse(c, "Invalid request payload", err.Error())
		return
	}

	consumer, err := req.ToEntity(primitive.NilObjectID)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid consumer payload", err.Error())
		return
	}

	if err := h.consumerService.CreateConsumer(c.Request.Context(), consumer); err != nil {
		switch {
		case errors.Is(err, services.ErrConsumerRepositoryUnavailable):
			response.NewInternalServerErrorResponse(c, "Consumer data unavailable", err.Error())
		case errors.Is(err, entities.ErrConsumerPersonalDataRequired),
			errors.Is(err, entities.ErrConsumerContactDataRequired),
			errors.Is(err, entities.ErrConsumerCreditProfileRequired),
			errors.Is(err, entities.ErrConsumerPrimaryAddressRequired):
			response.NewBadRequestResponse(c, "Consumer validation failed", err.Error())
		default:
			response.NewBadRequestResponse(c, "Unable to create consumer", err.Error())
		}
		return
	}

	response.NewCreatedResponse(c, "Consumer created successfully", dto.NewConsumerResponse(consumer))
}

func (h *ConsumerHandler) UpdateConsumer(c *gin.Context) {
	if h == nil || h.consumerService == nil {
		response.NewInternalServerErrorResponse(c, "Consumer service unavailable", "consumer service not configured")
		return
	}

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid consumer ID", err.Error())
		return
	}

	var req dto.ConsumerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewBadRequestResponse(c, "Invalid request payload", err.Error())
		return
	}

	consumer, err := req.ToEntity(id)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid consumer payload", err.Error())
		return
	}

	if err := h.consumerService.UpdateConsumer(c.Request.Context(), consumer); err != nil {
		switch {
		case errors.Is(err, services.ErrConsumerRepositoryUnavailable):
			response.NewInternalServerErrorResponse(c, "Consumer data unavailable", err.Error())
		case errors.Is(err, entities.ErrConsumerPersonalDataRequired),
			errors.Is(err, entities.ErrConsumerContactDataRequired),
			errors.Is(err, entities.ErrConsumerCreditProfileRequired),
			errors.Is(err, entities.ErrConsumerPrimaryAddressRequired):
			response.NewBadRequestResponse(c, "Consumer validation failed", err.Error())
		case errors.Is(err, services.ErrConsumerNotFound):
			response.NewNotFoundResponse(c, "Consumer not found", err.Error())
		default:
			response.NewBadRequestResponse(c, "Unable to update consumer", err.Error())
		}
		return
	}

	response.NewSuccessResponse(c, "Consumer updated successfully", dto.NewConsumerResponse(consumer))
}

func (h *ConsumerHandler) DeleteConsumer(c *gin.Context) {
	if h == nil || h.consumerService == nil {
		response.NewInternalServerErrorResponse(c, "Consumer service unavailable", "consumer service not configured")
		return
	}

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid consumer ID", err.Error())
		return
	}

	existingConsumer, err := h.consumerService.GetConsumerByID(c.Request.Context(), id)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to retrieve consumer", err.Error())
		return
	}

	if existingConsumer == nil {
		response.NewNotFoundResponse(c, "Consumer not found", "Consumer with the given ID does not exist")
		return
	}

	if err := h.consumerService.DeleteConsumer(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, services.ErrConsumerRepositoryUnavailable):
			response.NewInternalServerErrorResponse(c, "Consumer data unavailable", err.Error())
		default:
			response.NewInternalServerErrorResponse(c, "Failed to delete consumer", err.Error())
		}
		return
	}

	response.NewDeleteSuccessResponse(c, "Consumer", id.Hex())
}

func (h *ConsumerHandler) ListConsumers(c *gin.Context) {
	if h == nil || h.consumerService == nil {
		response.NewInternalServerErrorResponse(c, "Consumer service unavailable", "consumer service not configured")
		return
	}

	consumers, err := h.consumerService.ListConsumers(c.Request.Context(), map[string]interface{}{})
	if err != nil {
		switch {
		case errors.Is(err, services.ErrConsumerRepositoryUnavailable):
			response.NewInternalServerErrorResponse(c, "Consumer data unavailable", err.Error())
		default:
			response.NewInternalServerErrorResponse(c, "Failed to retrieve consumers", err.Error())
		}
		return
	}

	response.NewSuccessResponse(c, "Consumers retrieved successfully", dto.NewConsumerResponseList(consumers))
}

func (h *ConsumerHandler) ContractProduct(c *gin.Context) {
	if h == nil || h.consumerService == nil {
		response.NewInternalServerErrorResponse(c, "Consumer service unavailable", "consumer service not configured")
		return
	}

	consumerID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid consumer ID", err.Error())
		return
	}

	productID, err := primitive.ObjectIDFromHex(c.Param("product_id"))
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid product ID", err.Error())
		return
	}

	if err := h.consumerService.ContractProduct(c.Request.Context(), consumerID, productID); err != nil {
		switch {
		case errors.Is(err, services.ErrConsumerNotFound):
			response.NewNotFoundResponse(c, "Consumer not found", err.Error())
		case errors.Is(err, services.ErrProductNotFound):
			response.NewNotFoundResponse(c, "Product not found", err.Error())
		case errors.Is(err, entities.ErrConsumerProductAlreadyLinked):
			response.NewConflictResponse(c, "Product already contracted", err.Error())
		case errors.Is(err, services.ErrConsumerRepositoryUnavailable),
			errors.Is(err, services.ErrProductRepositoryUnavailable):
			response.NewInternalServerErrorResponse(c, "Operation unavailable", err.Error())
		default:
			response.NewInternalServerErrorResponse(c, "Failed to contract product", err.Error())
		}
		return
	}

	updatedConsumer, err := h.consumerService.GetConsumerByID(c.Request.Context(), consumerID)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to fetch updated consumer", err.Error())
		return
	}

	response.NewSuccessResponse(c, "Product contracted successfully", dto.NewConsumerResponse(updatedConsumer))
}

func (h *ConsumerHandler) RemoveProduct(c *gin.Context) {
	if h == nil || h.consumerService == nil {
		response.NewInternalServerErrorResponse(c, "Consumer service unavailable", "consumer service not configured")
		return
	}

	consumerID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid consumer ID", err.Error())
		return
	}

	productID, err := primitive.ObjectIDFromHex(c.Param("product_id"))
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid product ID", err.Error())
		return
	}

	if err := h.consumerService.RemoveContractedProduct(c.Request.Context(), consumerID, productID); err != nil {
		switch {
		case errors.Is(err, services.ErrConsumerNotFound):
			response.NewNotFoundResponse(c, "Consumer not found", err.Error())
		case errors.Is(err, entities.ErrConsumerProductNotLinked):
			response.NewNotFoundResponse(c, "Relationship not found", err.Error())
		case errors.Is(err, services.ErrConsumerRepositoryUnavailable):
			response.NewInternalServerErrorResponse(c, "Consumer data unavailable", err.Error())
		default:
			response.NewInternalServerErrorResponse(c, "Failed to remove contracted product", err.Error())
		}
		return
	}

	updatedConsumer, err := h.consumerService.GetConsumerByID(c.Request.Context(), consumerID)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to fetch updated consumer", err.Error())
		return
	}

	response.NewSuccessResponse(c, "Product contract removed successfully", dto.NewConsumerResponse(updatedConsumer))
}
