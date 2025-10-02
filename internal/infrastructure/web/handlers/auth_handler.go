package handlers

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"katseye/internal/domain/entities"
	"katseye/internal/domain/services"
	"katseye/internal/infrastructure/web/dto"
	"katseye/internal/infrastructure/web/response"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const defaultTokenTTL = 24 * time.Hour
const claimsContextKey = "jwt_claims"
const rawTokenContextKey = "jwt_raw_token"

type AuthHandler struct {
	authService  *services.AuthService
	tokenService *services.TokenService
	secret       []byte
	tokenTTL     time.Duration
}

func NewAuthHandler(service *services.AuthService, tokenService *services.TokenService, secret string) *AuthHandler {
	secret = strings.TrimSpace(secret)
	if service == nil || secret == "" {
		return nil
	}

	return &AuthHandler{
		authService:  service,
		tokenService: tokenService,
		secret:       []byte(secret),
		tokenTTL:     defaultTokenTTL,
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token       string   `json:"token"`
	ExpiresIn   int64    `json:"expires_in"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type createUserRequest struct {
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Active      *bool    `json:"active,omitempty"`
	Role        string   `json:"role,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
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
		Token:       token,
		ExpiresIn:   int64(h.tokenTTL.Seconds()),
		Role:        user.Role.String(),
		Permissions: user.Permissions,
	})
}

func (h *AuthHandler) CreateUser(c *gin.Context) {
	if h == nil || h.authService == nil {
		response.NewInternalServerErrorResponse(c, "Authentication service unavailable", "handler not configured")
		return
	}

	if _, ok := h.authorizeUserManagement(c); !ok {
		return
	}

	var req createUserRequest
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

	active := true
	if req.Active != nil {
		active = *req.Active
	}
	role, err := entities.ParseRole(req.Role)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid role", err.Error())
		return
	}

	user, err := h.authService.CreateUser(c.Request.Context(), email, password, active, role, req.Permissions)
	if err != nil {
		switch err {
		case services.ErrInvalidUserData:
			response.NewBadRequestResponse(c, "Invalid user data", err.Error())
		case services.ErrInvalidRole:
			response.NewBadRequestResponse(c, "Invalid role", err.Error())
		case services.ErrUserAlreadyExists:
			response.NewConflictResponse(c, "User already exists", err.Error())
		default:
			response.NewInternalServerErrorResponse(c, "Failed to create user", err.Error())
		}
		return
	}

	response.NewCreatedResponse(c, "User created successfully", dto.NewUserResponse(user))
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	if h == nil || h.authService == nil {
		response.NewInternalServerErrorResponse(c, "Authentication service unavailable", "handler not configured")
		return
	}

	actor, ok := h.authorizeUserManagement(c)
	if !ok {
		return
	}

	userID := strings.TrimSpace(c.Param("id"))
	if userID == "" {
		response.NewBadRequestResponse(c, "Invalid user ID", "user identifier is required")
		return
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid user ID", err.Error())
		return
	}

	if actor.ID == id {
		response.NewForbiddenResponse(c, "Cannot delete current user", "self-deletion is not permitted")
		return
	}

	if err := h.authService.DeleteUser(c.Request.Context(), id); err != nil {
		switch err {
		case services.ErrInvalidUserData:
			response.NewBadRequestResponse(c, "Invalid user data", err.Error())
		case services.ErrUserNotFound:
			response.NewNotFoundResponse(c, "User not found", err.Error())
		default:
			response.NewInternalServerErrorResponse(c, "Failed to delete user", err.Error())
		}
		return
	}

	response.NewDeleteSuccessResponse(c, "User", id.Hex())
}

func (h *AuthHandler) Logout(c *gin.Context) {
	if h == nil {
		response.NewUnauthorizedResponse(c, "Unauthorized", "authentication handler not configured")
		return
	}

	claimsValue, _ := c.Get(claimsContextKey)
	rawTokenValue, _ := c.Get(rawTokenContextKey)

	claims, _ := claimsValue.(jwt.MapClaims)
	rawToken, _ := rawTokenValue.(string)

	message := "Token invalidated on client side"

	if h.tokenService != nil && rawToken != "" && claims != nil {
		if expiresAt, ok := extractExpiration(claims); ok {
			if err := h.tokenService.RevokeToken(c.Request.Context(), rawToken, expiresAt); err != nil {
				response.NewInternalServerErrorResponse(c, "Failed to revoke token", err.Error())
				return
			}
			message = "Token revoked"
		}
	}

	response.NewSuccessResponse(c, "Logout successful", gin.H{
		"message": message,
	})
}

func (h *AuthHandler) currentUser(c *gin.Context) (*entities.User, error) {
	if h == nil || h.authService == nil {
		return nil, errors.New("authentication handler not configured")
	}

	claimsValue, exists := c.Get(claimsContextKey)
	if !exists {
		return nil, services.ErrInvalidCredentials
	}

	claims, ok := claimsValue.(jwt.MapClaims)
	if !ok {
		return nil, services.ErrInvalidCredentials
	}

	sub, ok := claims["sub"].(string)
	if !ok || strings.TrimSpace(sub) == "" {
		return nil, services.ErrInvalidCredentials
	}

	userID, err := primitive.ObjectIDFromHex(sub)
	if err != nil {
		return nil, services.ErrInvalidCredentials
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *AuthHandler) authorizeUserManagement(c *gin.Context) (*entities.User, bool) {
	user, err := h.currentUser(c)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			response.NewUnauthorizedResponse(c, "Unauthorized", "user not found")
		case services.ErrInvalidUserData, services.ErrInvalidCredentials:
			response.NewUnauthorizedResponse(c, "Unauthorized", "invalid authentication context")
		default:
			response.NewUnauthorizedResponse(c, "Unauthorized", "failed to resolve authenticated user")
		}
		return nil, false
	}

	if user.HasAnyRole(entities.RoleAdmin, entities.RoleManager) || user.HasPermission(entities.PermissionManageUsers) {
		return user, true
	}

	response.NewForbiddenResponse(c, "Forbidden", "insufficient permissions")
	return nil, false
}

func extractExpiration(claims jwt.MapClaims) (time.Time, bool) {
	if claims == nil {
		return time.Time{}, false
	}

	exp, err := claims.GetExpirationTime()
	if err != nil || exp == nil {
		return time.Time{}, false
	}

	return exp.Time, true
}

func (h *AuthHandler) generateToken(user *entities.User) (string, error) {
	if h == nil || user == nil {
		return "", services.ErrInvalidCredentials
	}

	if user.ID.IsZero() {
		return "", errors.New("user identifier is not set")
	}

	claims := jwt.MapClaims{
		"sub":         user.ID.Hex(),
		"email":       user.Email,
		"exp":         time.Now().Add(h.tokenTTL).Unix(),
		"iat":         time.Now().Unix(),
		"role":        user.Role.String(),
		"permissions": user.Permissions,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.secret)
}
