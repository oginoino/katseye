package valueobjects

import (
	"errors"
	"strings"
)

type PartnerType string

const (
	PartnerTypeBank               PartnerType = "bank"
	PartnerTypeCooperative        PartnerType = "cooperative"
	PartnerTypeFintechSCD         PartnerType = "fintech_scd"
	PartnerTypeFinanceira         PartnerType = "financeira"
	PartnerTypePaymentInstitution PartnerType = "payment_institution"
	PartnerTypeSavingsBank        PartnerType = "savings_bank"
	PartnerTypeDevelopmentBank    PartnerType = "development_bank"
)

var (
	ErrInvalidPartnerType = errors.New("invalid partner type")

	validPartnerTypes = map[PartnerType]bool{
		PartnerTypeBank:               true,
		PartnerTypeCooperative:        true,
		PartnerTypeFintechSCD:         true,
		PartnerTypeFinanceira:         true,
		PartnerTypePaymentInstitution: true,
		PartnerTypeSavingsBank:        true,
		PartnerTypeDevelopmentBank:    true,
	}
)

// NewPartnerType creates a validated PartnerType
func NewPartnerType(value string) (PartnerType, error) {
	pt := PartnerType(strings.ToLower(value))
	if err := pt.Validate(); err != nil {
		return "", err
	}
	return pt, nil
}

// Validate checks if the partner type is valid
func (pt PartnerType) Validate() error {
	if !validPartnerTypes[pt] {
		return ErrInvalidPartnerType
	}
	return nil
}

// String returns the string representation
func (pt PartnerType) String() string {
	return string(pt)
}

// IsRegulatedInstitution returns true if the partner type is a regulated financial institution
func (pt PartnerType) IsRegulatedInstitution() bool {
	regulatedTypes := []PartnerType{
		PartnerTypeBank,
		PartnerTypeCooperative,
		PartnerTypeSavingsBank,
		PartnerTypeDevelopmentBank,
	}

	for _, regulated := range regulatedTypes {
		if pt == regulated {
			return true
		}
	}
	return false
}

// GetRegulationRequirements returns the regulation requirements for this partner type
func (pt PartnerType) GetRegulationRequirements() []string {
	switch pt {
	case PartnerTypeBank:
		return []string{"bank_code", "bacen_authorization", "capital_requirement"}
	case PartnerTypeCooperative:
		return []string{"cooperative_registry", "ocb_registration", "audit_report"}
	case PartnerTypeFintechSCD:
		return []string{"scd_license", "bcb_authorization", "compliance_report"}
	case PartnerTypeFinanceira:
		return []string{"scfi_license", "financial_statements", "risk_rating"}
	case PartnerTypePaymentInstitution:
		return []string{"payment_license", "aml_compliance", "security_audit"}
	default:
		return []string{}
	}
}
