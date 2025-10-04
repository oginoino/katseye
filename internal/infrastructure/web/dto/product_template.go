package dto

import "katseye/internal/domain/services"

type ProductTemplateResponse struct {
	ProductType string      `json:"product_type"`
	DisplayName string      `json:"display_name"`
	Category    string      `json:"category"`
	Attributes  interface{} `json:"attributes"`
}

func NewProductTemplateResponse(template services.ProductTemplate) ProductTemplateResponse {
	return ProductTemplateResponse{
		ProductType: template.Type.String(),
		DisplayName: template.DisplayName,
		Category:    template.Category,
		Attributes:  template.Attributes,
	}
}

func NewProductTemplateListResponse(templates []services.ProductTemplate) []ProductTemplateResponse {
	if len(templates) == 0 {
		return nil
	}

	responses := make([]ProductTemplateResponse, 0, len(templates))
	for _, template := range templates {
		responses = append(responses, NewProductTemplateResponse(template))
	}
	return responses
}
