package handlers

import (
	"katseye/internal/domain/services"
	"katseye/internal/infrastructure/web/dto"
	"katseye/internal/infrastructure/web/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddressHandler struct {
	addressService *services.AddressService
}

func NewAddressHandler(addressService *services.AddressService) *AddressHandler {
	return &AddressHandler{
		addressService: addressService,
	}
}

func (h *AddressHandler) GetAddress(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid address ID", err.Error())
		return
	}

	address, err := h.addressService.GetAddressByID(c.Request.Context(), id)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to retrieve address", err.Error())
		return
	}

	if address == nil {
		response.NewNotFoundResponse(c, "Address not found", "Address with the given ID does not exist")
		return
	}

	response.NewSuccessResponse(c, "Address retrieved successfully", dto.NewAddressResponse(address))
}

func (h *AddressHandler) CreateAddress(c *gin.Context) {
	var req dto.AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewBadRequestResponse(c, "Invalid request payload", err.Error())
		return
	}

	address, err := req.ToEntity(primitive.NilObjectID)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid address payload", err.Error())
		return
	}

	if err := h.addressService.CreateAddress(c.Request.Context(), address); err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to create address", err.Error())
		return
	}

	response.NewCreatedResponse(c, "Address created successfully", dto.NewAddressResponse(address))
}

func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid address ID", err.Error())
		return
	}

	var req dto.AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewBadRequestResponse(c, "Invalid request payload", err.Error())
		return
	}

	address, err := req.ToEntity(id)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid address payload", err.Error())
		return
	}

	if err := h.addressService.UpdateAddress(c.Request.Context(), address); err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to update address", err.Error())
		return
	}

	response.NewSuccessResponse(c, "Address updated successfully", dto.NewAddressResponse(address))
}

func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid address ID", err.Error())
		return
	}

	existingAddress, err := h.addressService.GetAddressByID(c.Request.Context(), id)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to retrieve address", err.Error())
		return
	}

	if existingAddress == nil {
		response.NewNotFoundResponse(c, "Address not found", "Address with the given ID does not exist")
		return
	}

	if err := h.addressService.DeleteAddress(c.Request.Context(), id); err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to delete address", err.Error())
		return
	}

	response.NewDeleteSuccessResponse(c, "Address", id.Hex())
}

func (h *AddressHandler) ListAddresses(c *gin.Context) {
	filter := make(map[string]interface{})

	addresses, err := h.addressService.ListAddresses(c.Request.Context(), filter)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to retrieve addresses", err.Error())
		return
	}

	response.NewSuccessResponse(c, "Addresses retrieved successfully", dto.NewAddressResponseList(addresses))
}
