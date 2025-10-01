package handlers

import (
	"katseye/internal/domain/entities"
	"katseye/internal/domain/services"
	"katseye/internal/infrastructure/web/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PartnerHandler struct {
	partnerService *services.PartnerService
}

func NewPartnerHandler(partnerService *services.PartnerService) *PartnerHandler {
	return &PartnerHandler{
		partnerService: partnerService,
	}
}

func (h *PartnerHandler) GetPartner(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid partner ID", err.Error())
		return
	}

	partner, err := h.partnerService.GetPartnerByID(c.Request.Context(), id)
	if err != nil {
		response.NewNotFoundResponse(c, "Partner not found", err.Error())
		return
	}

	if partner == nil {
		response.NewNotFoundResponse(c, "Partner not found", "Partner with the given ID does not exist")
		return
	}

	response.NewSuccessResponse(c, "Partner retrieved successfully", partner)
}

func (h *PartnerHandler) CreatePartner(c *gin.Context) {
	var partner entities.Partner
	if err := c.ShouldBindJSON(&partner); err != nil {
		response.NewBadRequestResponse(c, "Invalid request payload", err.Error())
		return
	}

	err := h.partnerService.CreatePartner(c.Request.Context(), &partner)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to create partner", err.Error())
		return
	}

	response.NewCreatedResponse(c, "Partner created successfully", partner)
}

func (h *PartnerHandler) UpdatePartner(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid partner ID", err.Error())
		return
	}

	var partner entities.Partner
	if err := c.ShouldBindJSON(&partner); err != nil {
		response.NewBadRequestResponse(c, "Invalid request payload", err.Error())
		return
	}
	partner.ID = id

	err = h.partnerService.UpdatePartner(c.Request.Context(), &partner)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to update partner", err.Error())
		return
	}

	response.NewSuccessResponse(c, "Partner updated successfully", partner)
}

func (h *PartnerHandler) DeletePartner(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid partner ID", err.Error())
		return
	}

	// Optional: Check if partner exists before deleting
	existingPartner, err := h.partnerService.GetPartnerByID(c.Request.Context(), id)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to retrieve partner", err.Error())
		return
	}

	if existingPartner == nil {
		response.NewNotFoundResponse(c, "Partner not found", "Partner with the given ID does not exist")
		return
	}

	err = h.partnerService.DeletePartner(c.Request.Context(), id)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to delete partner", err.Error())
		return
	}

	response.NewDeleteSuccessResponse(c, "Partner", id.Hex())
}

func (h *PartnerHandler) ListPartners(c *gin.Context) {
	// You can extend this to accept query parameters for filtering
	filter := make(map[string]interface{})

	partners, err := h.partnerService.ListPartners(c.Request.Context(), filter)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to list partners", err.Error())
		return
	}

	response.NewSuccessResponse(c, "Partners retrieved successfully", partners)
}
