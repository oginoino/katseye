package dto

import "katseye/internal/domain/entities"

// UserResponse representa o payload exposto para consumidores HTTP.
type UserResponse struct {
	ID          string   `json:"id"`
	Email       string   `json:"email"`
	Active      bool     `json:"active"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

// NewUserResponse converte a entidade de domínio em DTO.
func NewUserResponse(user *entities.User) UserResponse {
	if user == nil {
		return UserResponse{}
	}

	return UserResponse{
		ID:          user.ID.Hex(),
		Email:       user.Email,
		Active:      user.Active,
		Role:        user.Role.String(),
		Permissions: append([]string(nil), user.Permissions...),
	}
}
