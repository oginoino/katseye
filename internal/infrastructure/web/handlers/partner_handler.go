package handlers

import (
	"katseye/internal/domain/services"
	"katseye/internal/infrastructure/web/dto"
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
		response.NewInternalServerErrorResponse(c, "Failed to retrieve partner", err.Error())
		return
	}

	if partner == nil {
		response.NewNotFoundResponse(c, "Partner not found", "Partner with the given ID does not exist")
		return
	}

	response.NewSuccessResponse(c, "Partner retrieved successfully", dto.NewPartnerResponse(partner))
}

func (h *PartnerHandler) CreatePartner(c *gin.Context) {
	var req dto.PartnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewBadRequestResponse(c, "Invalid request payload", err.Error())
		return
	}

	partner, err := req.ToEntity(primitive.NilObjectID)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid partner payload", err.Error())
		return
	}

	if err := h.partnerService.CreatePartner(c.Request.Context(), partner); err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to create partner", err.Error())
		return
	}

	response.NewCreatedResponse(c, "Partner created successfully", dto.NewPartnerResponse(partner))
}

func (h *PartnerHandler) UpdatePartner(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid partner ID", err.Error())
		return
	}

	var req dto.PartnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewBadRequestResponse(c, "Invalid request payload", err.Error())
		return
	}

	partner, err := req.ToEntity(id)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid partner payload", err.Error())
		return
	}

	if err := h.partnerService.UpdatePartner(c.Request.Context(), partner); err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to update partner", err.Error())
		return
	}

	response.NewSuccessResponse(c, "Partner updated successfully", dto.NewPartnerResponse(partner))
}

func (h *PartnerHandler) DeletePartner(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid partner ID", err.Error())
		return
	}

	existingPartner, err := h.partnerService.GetPartnerByID(c.Request.Context(), id)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to retrieve partner", err.Error())
		return
	}

	if existingPartner == nil {
		response.NewNotFoundResponse(c, "Partner not found", "Partner with the given ID does not exist")
		return
	}

	if err := h.partnerService.DeletePartner(c.Request.Context(), id); err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to delete partner", err.Error())
		return
	}

	response.NewDeleteSuccessResponse(c, "Partner", id.Hex())
}

func (h *PartnerHandler) ListPartners(c *gin.Context) {
	filter := make(map[string]interface{})

	partners, err := h.partnerService.ListPartners(c.Request.Context(), filter)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to list partners", err.Error())
		return
	}

	response.NewSuccessResponse(c, "Partners retrieved successfully", dto.NewPartnerResponseList(partners))
}
