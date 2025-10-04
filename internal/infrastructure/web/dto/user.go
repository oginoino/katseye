package dto

import "katseye/internal/domain/entities"

// UserResponse representa o payload exposto para consumidores HTTP.
type UserResponse struct {
	ID          string   `json:"id"`
	Email       string   `json:"email"`
	Active      bool     `json:"active"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	ProfileType string   `json:"profile_type"`
	ProfileID   string   `json:"profile_reference_id,omitempty"`
}

// NewUserResponse converte a entidade de domínio em DTO.
func NewUserResponse(user *entities.User) UserResponse {
	if user == nil {
		return UserResponse{}
	}

	response := UserResponse{
		ID:          user.ID.Hex(),
		Email:       user.Email,
		Active:      user.Active,
		Role:        user.Role.String(),
		Permissions: append([]string(nil), user.Permissions...),
		ProfileType: user.ProfileType.String(),
	}

	if !user.ProfileID.IsZero() {
		response.ProfileID = user.ProfileID.Hex()
	}

	return response
}
