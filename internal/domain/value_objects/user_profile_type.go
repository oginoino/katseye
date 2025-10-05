package valueobjects

import "errors"

type UserProfileType string

const (
	UserProfileTypeServiceAccount UserProfileType = "service_account"
	UserProfileTypePartnerManager UserProfileType = "partner_manager"
	UserProfileTypeConsumer       UserProfileType = "consumer"
)

var (
	ErrInvalidUserProfileType = errors.New("invalid user profile type")
)

func (p UserProfileType) String() string {
	return string(p)
}

func (p UserProfileType) Validate() error {
	switch p {
	case UserProfileTypeServiceAccount, UserProfileTypePartnerManager, UserProfileTypeConsumer:
		return nil
	default:
		return ErrInvalidUserProfileType
	}
}
