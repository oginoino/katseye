package handlers

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"katseye/internal/domain/entities"
	"katseye/internal/domain/services"
	"katseye/internal/infrastructure/web/response"
)

const defaultTokenTTL = 24 * time.Hour

type AuthHandler struct {
	authService *services.AuthService
	secret      []byte
	tokenTTL    time.Duration
}

func NewAuthHandler(service *services.AuthService, secret string) *AuthHandler {
	secret = strings.TrimSpace(secret)
	if service == nil || secret == "" {
		return nil
	}

	return &AuthHandler{
		authService: service,
		secret:      []byte(secret),
		tokenTTL:    defaultTokenTTL,
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	if h == nil || h.authService == nil {
		response.NewInternalServerErrorResponse(c, "Authentication service unavailable", "handler not configured")
		return
	}

	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewBadRequestResponse(c, "Invalid request payload", err.Error())
		return
	}

	email := strings.TrimSpace(req.Email)
	password := strings.TrimSpace(req.Password)
	if email == "" || password == "" {
		response.NewBadRequestResponse(c, "Email and password are required", "missing credentials")
		return
	}

	ctx := c.Request.Context()
	user, err := h.authService.Authenticate(ctx, email, password)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			response.NewUnauthorizedResponse(c, "Invalid credentials", err.Error())
		case services.ErrInactiveAccount:
			response.NewForbiddenResponse(c, "Account is inactive", err.Error())
		default:
			response.NewInternalServerErrorResponse(c, "Failed to authenticate", err.Error())
		}
		return
	}

	token, err := h.generateToken(user)
	if err != nil {
		response.NewInternalServerErrorResponse(c, "Failed to generate token", err.Error())
		return
	}

	response.NewSuccessResponse(c, "Authentication successful", loginResponse{
		Token:     token,
		ExpiresIn: int64(h.tokenTTL.Seconds()),
	})
}

func (h *AuthHandler) generateToken(user *entities.User) (string, error) {
	if h == nil || user == nil {
		return "", services.ErrInvalidCredentials
	}

	if user.ID.IsZero() {
		return "", errors.New("user identifier is not set")
	}

	claims := jwt.MapClaims{
		"sub":   user.ID.Hex(),
		"email": user.Email,
		"exp":   time.Now().Add(h.tokenTTL).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.secret)
}
