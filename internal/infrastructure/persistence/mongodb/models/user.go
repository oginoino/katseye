package models

import (
	"katseye/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserDocument descreve como usuários são persistidos no MongoDB.
type UserDocument struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"password_hash"`
	Active       bool               `bson:"active"`
	Role         string             `bson:"role"`
	Permissions  []string           `bson:"permissions"`
	ProfileType  string             `bson:"profile_type"`
	ProfileID    primitive.ObjectID `bson:"profile_id,omitempty"`
}

// ToEntity converte o documento em entidade de domínio.
func (doc UserDocument) ToEntity() *entities.User {
	role, err := entities.ParseRole(doc.Role)
	if err != nil {
		role = entities.RoleUser
	}

	user := &entities.User{
		ID:           doc.ID,
		Email:        doc.Email,
		PasswordHash: doc.PasswordHash,
		Active:       doc.Active,
		Role:         role,
		Permissions:  append([]string(nil), doc.Permissions...),
		ProfileType:  entities.UserProfileType(doc.ProfileType),
		ProfileID:    doc.ProfileID,
	}
	user.Normalize()

	return user
}

// NewUserDocument converte uma entidade de domínio em documento persistido.
func NewUserDocument(user *entities.User) UserDocument {
	if user == nil {
		return UserDocument{}
	}

	normalized := *user
	normalized.Normalize()

	return UserDocument{
		ID:           normalized.ID,
		Email:        normalized.Email,
		PasswordHash: normalized.PasswordHash,
		Active:       normalized.Active,
		Role:         normalized.Role.String(),
		Permissions:  append([]string(nil), normalized.Permissions...),
		ProfileType:  normalized.ProfileType.String(),
		ProfileID:    normalized.ProfileID,
	}
}
