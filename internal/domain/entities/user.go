package entities

import (
	"errors"
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Role represents the authorization level of a user account.
type Role string

const (
	RoleAdmin   Role = "admin"
	RoleManager Role = "manager"
	RoleUser    Role = "user"

	PermissionManageUsers    = "users:manage"
	PermissionManageProducts = "products:manage"
	PermissionManagePartners = "partners:manage"

	PermissionEditUsers    = "users:edit"
	PermissionEditProducts = "products:edit"
	PermissionEditPartners = "partners:edit"

	PermissionViewUsers    = "users:view"
	PermissionViewProducts = "products:view"
	PermissionViewPartners = "partners:view"
)

// rolePermissions defines the base permissions for each role
var rolePermissions = map[Role][]string{
	RoleAdmin: {
		PermissionManageUsers,
		PermissionManageProducts,
		PermissionManagePartners,
		PermissionEditUsers,
		PermissionEditProducts,
		PermissionEditPartners,
		PermissionViewUsers,
		PermissionViewProducts,
		PermissionViewPartners,
	},
	RoleManager: {
		PermissionEditProducts,
		PermissionEditPartners,
		PermissionViewUsers,
		PermissionViewProducts,
		PermissionViewPartners,
	},
	RoleUser: {
		PermissionViewProducts,
	},
}

// UserProfileType represents the type of profile associated with given credentials.
type UserProfileType string

const (
	ProfileTypeServiceAccount UserProfileType = "service_account"
	ProfileTypePartnerManager UserProfileType = "partner_manager"
	ProfileTypeConsumer       UserProfileType = "consumer"
)

var (
	ErrInvalidPassword    = errors.New("invalid password")
	ErrEmptyPassword      = errors.New("password must not be empty")
	ErrInvalidRole        = errors.New("invalid role")
	ErrInvalidProfileType = errors.New("invalid profile type")
)

type User struct {
	ID           primitive.ObjectID
	Email        string
	PasswordHash string
	Active       bool
	Role         Role
	Permissions  []string // Custom permissions in addition to role-based permissions
	ProfileType  UserProfileType
	ProfileID    primitive.ObjectID
}

// Normalize prepares user fields for persistence/lookup.
func (u *User) Normalize() {
	if u == nil {
		return
	}
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	role := Role(strings.TrimSpace(strings.ToLower(u.Role.String())))
	if role == "" {
		role = RoleUser
	}
	if !IsValidRole(role) {
		role = RoleUser
	}
	u.Role = role

	// Normalize only custom permissions (role-based permissions are computed)
	u.Permissions = normalizePermissions(u.Permissions)

	profileType := UserProfileType(strings.TrimSpace(strings.ToLower(u.ProfileType.String())))
	if profileType == "" {
		profileType = ProfileTypeServiceAccount
	}
	if !IsValidProfileType(profileType) {
		profileType = ProfileTypeServiceAccount
	}
	u.ProfileType = profileType

	if u.ProfileType == ProfileTypeServiceAccount {
		u.ProfileID = primitive.NilObjectID
	}
}

// GetEffectivePermissions returns all permissions for this user, combining
// role-based permissions with custom permissions.
func (u *User) GetEffectivePermissions() []string {
	if u == nil {
		return nil
	}

	// Get base permissions for the role
	basePerms := rolePermissions[u.Role]
	if basePerms == nil {
		basePerms = []string{}
	}

	// Combine with custom permissions
	allPerms := make([]string, 0, len(basePerms)+len(u.Permissions))
	allPerms = append(allPerms, basePerms...)
	allPerms = append(allPerms, u.Permissions...)

	return normalizePermissions(allPerms)
}

// GetRolePermissions returns the base permissions for a given role.
func GetRolePermissions(role Role) []string {
	perms := rolePermissions[role]
	if perms == nil {
		return []string{}
	}
	// Return a copy to prevent external modification
	result := make([]string, len(perms))
	copy(result, perms)
	return result
}

// CheckPassword compares a clear-text password with the stored bcrypt hash.
func (u *User) CheckPassword(password string) error {
	if u == nil {
		return ErrInvalidPassword
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return ErrInvalidPassword
	}
	return nil
}

// SetPassword hashes and stores the provided clear-text password.
func (u *User) SetPassword(password string) error {
	if u == nil {
		return ErrInvalidPassword
	}
	password = strings.TrimSpace(password)
	if password == "" {
		return ErrEmptyPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// IsActive returns true when the account is enabled for authentication.
func (u *User) IsActive() bool {
	if u == nil {
		return false
	}
	return u.Active
}

// HasAnyRole returns true when the user possesses at least one of the provided roles.
func (u *User) HasAnyRole(roles ...Role) bool {
	if u == nil || len(roles) == 0 {
		return false
	}
	for _, role := range roles {
		if strings.EqualFold(role.String(), u.Role.String()) {
			return true
		}
	}
	return false
}

// HasPermission returns true when the user has the given permission (case insensitive).
// This checks both role-based and custom permissions.
func (u *User) HasPermission(permission string) bool {
	if u == nil {
		return false
	}
	permission = strings.TrimSpace(strings.ToLower(permission))
	if permission == "" {
		return false
	}

	// Check role-based permissions first
	rolePerms := rolePermissions[u.Role]
	for _, perm := range rolePerms {
		if perm == permission {
			return true
		}
	}

	// Check custom permissions
	for _, perm := range u.Permissions {
		if perm == permission {
			return true
		}
	}

	return false
}

// HasAllPermissions returns true when the user has all the given permissions.
func (u *User) HasAllPermissions(permissions ...string) bool {
	if u == nil {
		return false
	}
	for _, perm := range permissions {
		if !u.HasPermission(perm) {
			return false
		}
	}
	return true
}

// HasAnyPermission returns true when the user has at least one of the given permissions.
func (u *User) HasAnyPermission(permissions ...string) bool {
	if u == nil {
		return false
	}
	for _, perm := range permissions {
		if u.HasPermission(perm) {
			return true
		}
	}
	return false
}

// String returns the string representation of the role.
func (r Role) String() string {
	return string(r)
}

// IsValidRole reports whether the provided role belongs to the list of supported roles.
func IsValidRole(role Role) bool {
	switch Role(strings.TrimSpace(strings.ToLower(role.String()))) {
	case RoleAdmin, RoleManager, RoleUser:
		return true
	default:
		return false
	}
}

// ParseRole converts the provided string into a Role, validating the value.
func ParseRole(role string) (Role, error) {
	candidate := Role(strings.TrimSpace(strings.ToLower(role)))
	if candidate == "" {
		return RoleUser, nil
	}
	if !IsValidRole(candidate) {
		return "", ErrInvalidRole
	}
	return candidate, nil
}

// String returns the profile type as string.
func (t UserProfileType) String() string {
	return string(t)
}

// IsValidProfileType returns true when the provided profile type is recognised.
func IsValidProfileType(profileType UserProfileType) bool {
	switch UserProfileType(strings.TrimSpace(strings.ToLower(profileType.String()))) {
	case ProfileTypeServiceAccount, ProfileTypePartnerManager, ProfileTypeConsumer:
		return true
	default:
		return false
	}
}

// ParseProfileType normalizes and validates the profile type string.
func ParseProfileType(profileType string) (UserProfileType, error) {
	candidate := UserProfileType(strings.TrimSpace(strings.ToLower(profileType)))
	if candidate == "" {
		return ProfileTypeServiceAccount, nil
	}
	if !IsValidProfileType(candidate) {
		return "", ErrInvalidProfileType
	}
	return candidate, nil
}

func normalizePermissions(perms []string) []string {
	if len(perms) == 0 {
		return nil
	}

	set := make(map[string]struct{}, len(perms))
	for _, perm := range perms {
		perm = strings.TrimSpace(strings.ToLower(perm))
		if perm == "" {
			continue
		}
		set[perm] = struct{}{}
	}

	if len(set) == 0 {
		return nil
	}

	normalized := make([]string, 0, len(set))
	for perm := range set {
		normalized = append(normalized, perm)
	}
	sort.Strings(normalized)
	return normalized
}
