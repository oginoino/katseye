package handlers

import (
	"context"
	"errors"
	"fmt"
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
	authService     *services.AuthService
	tokenService    *services.TokenService
	partnerService  *services.PartnerService
	consumerService *services.ConsumerService
	secret          []byte
	tokenTTL        time.Duration
}

func NewAuthHandler(service *services.AuthService, tokenService *services.TokenService, partnerService *services.PartnerService, consumerService *services.ConsumerService, secret string) *AuthHandler {
	secret = strings.TrimSpace(secret)
	if service == nil || secret == "" {
		return nil
	}

	return &AuthHandler{
		authService:     service,
		tokenService:    tokenService,
		partnerService:  partnerService,
		consumerService: consumerService,
		secret:          []byte(secret),
		tokenTTL:        defaultTokenTTL,
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
	ProfileType string   `json:"profile_type"`
	ProfileID   string   `json:"profile_reference_id,omitempty"`
}

type createUserRequest struct {
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Active      *bool    `json:"active,omitempty"`
	Role        string   `json:"role,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	ProfileType string   `json:"profile_type,omitempty"`
	ProfileID   string   `json:"profile_reference_id,omitempty"`
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

	resp := loginResponse{
		Token:       token,
		ExpiresIn:   int64(h.tokenTTL.Seconds()),
		Role:        user.Role.String(),
		Permissions: user.Permissions,
		ProfileType: user.ProfileType.String(),
	}
	if !user.ProfileID.IsZero() {
		resp.ProfileID = user.ProfileID.Hex()
	}

	response.NewSuccessResponse(c, "Authentication successful", resp)
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

	profileType, err := entities.ParseProfileType(req.ProfileType)
	if err != nil {
		response.NewBadRequestResponse(c, "Invalid profile type", err.Error())
		return
	}

	var profileID primitive.ObjectID
	if profileType != entities.ProfileTypeServiceAccount {
		trimmed := strings.TrimSpace(req.ProfileID)
		if trimmed == "" {
			response.NewBadRequestResponse(c, "Profile reference is required", "profile_reference_id is required for the selected profile type")
			return
		}
		profileID, err = primitive.ObjectIDFromHex(trimmed)
		if err != nil {
			response.NewBadRequestResponse(c, "Invalid profile reference", err.Error())
			return
		}
	}

	ctx := c.Request.Context()

	switch profileType {
	case entities.ProfileTypePartnerManager:
		if h.partnerService == nil {
			response.NewInternalServerErrorResponse(c, "Partner service unavailable", "partner service not configured")
			return
		}
		partner, err := h.partnerService.GetPartnerByID(ctx, profileID)
		if err != nil {
			response.NewInternalServerErrorResponse(c, "Failed to validate partner", err.Error())
			return
		}
		if partner == nil {
			response.NewNotFoundResponse(c, "Partner not found", "linked partner does not exist")
			return
		}
	case entities.ProfileTypeConsumer:
		if h.consumerService == nil {
			response.NewInternalServerErrorResponse(c, "Consumer service unavailable", "consumer service not configured")
			return
		}
		consumer, err := h.consumerService.GetConsumerByID(ctx, profileID)
		if err != nil {
			response.NewInternalServerErrorResponse(c, "Failed to validate consumer", err.Error())
			return
		}
		if consumer == nil {
			response.NewNotFoundResponse(c, "Consumer not found", "linked consumer does not exist")
			return
		}
		if consumer.HasLinkedUser() {
			response.NewConflictResponse(c, "Consumer already linked", "the consumer already has an authentication profile")
			return
		}
	}

	user, err := h.authService.CreateUser(ctx, email, password, active, role, req.Permissions, profileType, profileID)
	if err != nil {
		switch err {
		case services.ErrInvalidUserData:
			response.NewBadRequestResponse(c, "Invalid user data", err.Error())
		case services.ErrInvalidRole:
			response.NewBadRequestResponse(c, "Invalid role", err.Error())
		case services.ErrInvalidProfileType:
			response.NewBadRequestResponse(c, "Invalid profile type", err.Error())
		case services.ErrUserAlreadyExists:
			response.NewConflictResponse(c, "User already exists", err.Error())
		default:
			response.NewInternalServerErrorResponse(c, "Failed to create user", err.Error())
		}
		return
	}

	switch profileType {
	case entities.ProfileTypePartnerManager:
		if err := h.partnerService.AssignManagerProfile(ctx, profileID, user.ID); err != nil {
			h.handleProfileLinkError(c, user.ID, err)
			return
		}
	case entities.ProfileTypeConsumer:
		if err := h.consumerService.AttachUserProfile(ctx, profileID, user.ID); err != nil {
			h.handleProfileLinkError(c, user.ID, err)
			return
		}
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

	ctx := c.Request.Context()
	target, err := h.authService.GetUserByID(ctx, id)
	if err != nil {
		switch err {
		case services.ErrInvalidUserData:
			response.NewBadRequestResponse(c, "Invalid user data", err.Error())
		case services.ErrUserNotFound:
			response.NewNotFoundResponse(c, "User not found", err.Error())
		default:
			response.NewInternalServerErrorResponse(c, "Failed to retrieve user", err.Error())
		}
		return
	}
	if target == nil {
		response.NewNotFoundResponse(c, "User not found", "user does not exist")
		return
	}

	rollback := func(context.Context) error { return nil }

	switch target.ProfileType {
	case entities.ProfileTypePartnerManager:
		if err := h.preparePartnerManagerRemoval(ctx, target); err != nil {
			h.respondProfileRemovalError(c, err)
			return
		}
		rollback = func(ctx context.Context) error {
			return h.reassignPartnerManager(ctx, target)
		}
	case entities.ProfileTypeConsumer:
		if err := h.prepareConsumerDetachment(ctx, target); err != nil {
			h.respondProfileRemovalError(c, err)
			return
		}
		rollback = func(ctx context.Context) error {
			return h.consumerService.AttachUserProfile(ctx, target.ProfileID, target.ID)
		}
	}

	if err := h.authService.DeleteUser(ctx, id); err != nil {
		if revertErr := rollback(ctx); revertErr != nil {
			err = fmt.Errorf("%w; rollback failed: %v", err, revertErr)
		}
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
		"sub":          user.ID.Hex(),
		"email":        user.Email,
		"exp":          time.Now().Add(h.tokenTTL).Unix(),
		"iat":          time.Now().Unix(),
		"role":         user.Role.String(),
		"permissions":  user.Permissions,
		"profile_type": user.ProfileType.String(),
	}
	if !user.ProfileID.IsZero() {
		claims["profile_reference_id"] = user.ProfileID.Hex()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.secret)
}

func (h *AuthHandler) handleProfileLinkError(c *gin.Context, userID primitive.ObjectID, linkErr error) {
	if linkErr == nil {
		return
	}

	if rollbackErr := h.rollbackUser(c.Request.Context(), userID); rollbackErr != nil {
		linkErr = fmt.Errorf("%w; rollback failed: %v", linkErr, rollbackErr)
	}

	switch {
	case errors.Is(linkErr, services.ErrPartnerNotFound):
		response.NewNotFoundResponse(c, "Partner not found", linkErr.Error())
	case errors.Is(linkErr, services.ErrPartnerManagerAlreadyLinked):
		response.NewConflictResponse(c, "Manager already linked to partner", linkErr.Error())
	case errors.Is(linkErr, services.ErrPartnerRepositoryUnavailable):
		response.NewInternalServerErrorResponse(c, "Partner service unavailable", linkErr.Error())
	case errors.Is(linkErr, services.ErrConsumerNotFound):
		response.NewNotFoundResponse(c, "Consumer not found", linkErr.Error())
	case errors.Is(linkErr, services.ErrConsumerUserAlreadyLinked):
		response.NewConflictResponse(c, "Consumer already linked", linkErr.Error())
	case errors.Is(linkErr, services.ErrConsumerRepositoryUnavailable):
		response.NewInternalServerErrorResponse(c, "Consumer service unavailable", linkErr.Error())
	default:
		response.NewInternalServerErrorResponse(c, "Failed to link user profile", linkErr.Error())
	}
}

func (h *AuthHandler) rollbackUser(ctx context.Context, userID primitive.ObjectID) error {
	if h == nil || h.authService == nil || userID.IsZero() {
		return nil
	}

	if err := h.authService.DeleteUser(ctx, userID); err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil
		}
		return err
	}

	return nil
}

func (h *AuthHandler) preparePartnerManagerRemoval(ctx context.Context, user *entities.User) error {
	if h == nil || h.partnerService == nil {
		return services.ErrPartnerRepositoryUnavailable
	}
	if user == nil || user.ProfileID.IsZero() {
		return nil
	}

	return h.partnerService.RemoveManagerProfile(ctx, user.ProfileID, user.ID)
}

func (h *AuthHandler) reassignPartnerManager(ctx context.Context, user *entities.User) error {
	if h == nil || h.partnerService == nil || user == nil || user.ProfileID.IsZero() {
		return nil
	}

	return h.partnerService.AssignManagerProfile(ctx, user.ProfileID, user.ID)
}

func (h *AuthHandler) prepareConsumerDetachment(ctx context.Context, user *entities.User) error {
	if h == nil || h.consumerService == nil {
		return services.ErrConsumerRepositoryUnavailable
	}
	if user == nil || user.ProfileID.IsZero() {
		return nil
	}

	return h.consumerService.DetachUserProfile(ctx, user.ProfileID)
}

func (h *AuthHandler) respondProfileRemovalError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrPartnerManagerRequired):
		response.NewConflictResponse(c, "Partner must retain a manager", err.Error())
	case errors.Is(err, services.ErrPartnerManagerNotLinked):
		response.NewNotFoundResponse(c, "Manager not linked to partner", err.Error())
	case errors.Is(err, services.ErrPartnerNotFound):
		response.NewNotFoundResponse(c, "Partner not found", err.Error())
	case errors.Is(err, services.ErrPartnerRepositoryUnavailable):
		response.NewInternalServerErrorResponse(c, "Partner service unavailable", err.Error())
	case errors.Is(err, services.ErrConsumerUserNotLinked):
		response.NewNotFoundResponse(c, "Consumer user link not found", err.Error())
	case errors.Is(err, services.ErrConsumerNotFound):
		response.NewNotFoundResponse(c, "Consumer not found", err.Error())
	case errors.Is(err, services.ErrConsumerRepositoryUnavailable):
		response.NewInternalServerErrorResponse(c, "Consumer service unavailable", err.Error())
	default:
		response.NewInternalServerErrorResponse(c, "Failed to update linked profile", err.Error())
	}
}
