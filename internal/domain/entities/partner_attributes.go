package entities

import (
	"errors"

	valueObjects "katseye/internal/domain/value_objects"
)

// BasePartnerAttributes contains common attributes for all partner types
type BasePartnerAttributes struct {
	RegulatoryLicense  string  `json:"regulatory_license,omitempty" bson:"regulatory_license,omitempty"`
	MinimumLoanAmount  float64 `json:"minimum_loan_amount,omitempty" bson:"minimum_loan_amount,omitempty"`
	MaximumLoanAmount  float64 `json:"maximum_loan_amount,omitempty" bson:"maximum_loan_amount,omitempty"`
	InterestRateRange  string  `json:"interest_rate_range,omitempty" bson:"interest_rate_range,omitempty"`
	ProcessingTimeDays int     `json:"processing_time_days,omitempty" bson:"processing_time_days,omitempty"`
}

// BankPartnerAttributes contains bank-specific attributes
type BankPartnerAttributes struct {
	BasePartnerAttributes `bson:",inline"`
	BankCode              string `json:"bank_code" bson:"bank_code"`
	BACENAuthorization    string `json:"bacen_authorization,omitempty" bson:"bacen_authorization,omitempty"`
	CapitalRequirement    string `json:"capital_requirement,omitempty" bson:"capital_requirement,omitempty"`
}

// CooperativePartnerAttributes contains cooperative-specific attributes
type CooperativePartnerAttributes struct {
	BasePartnerAttributes `bson:",inline"`
	CooperativeRegistry   string `json:"cooperative_registry" bson:"cooperative_registry"`
	OCBRegistration       string `json:"ocb_registration,omitempty" bson:"ocb_registration,omitempty"`
	AuditReportRequired   bool   `json:"audit_report_required" bson:"audit_report_required"`
}

// FintechSCDAttributes contains fintech SCD-specific attributes
type FintechSCDAttributes struct {
	BasePartnerAttributes `bson:",inline"`
	SCDLicense            string `json:"scd_license" bson:"scd_license"`
	BCBAuthorization      string `json:"bcb_authorization" bson:"bcb_authorization"`
	ComplianceReport      string `json:"compliance_report,omitempty" bson:"compliance_report,omitempty"`
}

// FinanceiraAttributes contains financeira-specific attributes
type FinanceiraAttributes struct {
	BasePartnerAttributes `bson:",inline"`
	SCFILicense           string  `json:"scfi_license" bson:"scfi_license"`
	RiskRating            string  `json:"risk_rating,omitempty" bson:"risk_rating,omitempty"`
	CapitalAdequacy       float64 `json:"capital_adequacy,omitempty" bson:"capital_adequacy,omitempty"`
}

// PaymentInstitutionAttributes contains payment institution-specific attributes
type PaymentInstitutionAttributes struct {
	BasePartnerAttributes `bson:",inline"`
	PaymentLicense        string `json:"payment_license" bson:"payment_license"`
	AMLCompliance         bool   `json:"aml_compliance" bson:"aml_compliance"`
	SecurityAudit         string `json:"security_audit,omitempty" bson:"security_audit,omitempty"`
}

// SavingsBankAttributes contains savings bank-specific attributes
type SavingsBankAttributes struct {
	BasePartnerAttributes `bson:",inline"`
	SavingsBankCode       string `json:"savings_bank_code" bson:"savings_bank_code"`
	PublicEntity          bool   `json:"public_entity" bson:"public_entity"`
}

// DevelopmentBankAttributes contains development bank-specific attributes
type DevelopmentBankAttributes struct {
	BasePartnerAttributes `bson:",inline"`
	DevelopmentBankCode   string   `json:"development_bank_code" bson:"development_bank_code"`
	TargetSectors         []string `json:"target_sectors" bson:"target_sectors"`
	GovernmentBacked      bool     `json:"government_backed" bson:"government_backed"`
}

// PartnerAttributes is a union type that can hold any partner-specific attributes
type PartnerAttributes struct {
	Bank               *BankPartnerAttributes        `json:"bank,omitempty" bson:"bank,omitempty"`
	Cooperative        *CooperativePartnerAttributes `json:"cooperative,omitempty" bson:"cooperative,omitempty"`
	FintechSCD         *FintechSCDAttributes         `json:"fintech_scd,omitempty" bson:"fintech_scd,omitempty"`
	Financeira         *FinanceiraAttributes         `json:"financeira,omitempty" bson:"financeira,omitempty"`
	PaymentInstitution *PaymentInstitutionAttributes `json:"payment_institution,omitempty" bson:"payment_institution,omitempty"`
	SavingsBank        *SavingsBankAttributes        `json:"savings_bank,omitempty" bson:"savings_bank,omitempty"`
	DevelopmentBank    *DevelopmentBankAttributes    `json:"development_bank,omitempty" bson:"development_bank,omitempty"`
}

// GetAttributesForPartnerType returns the appropriate attributes struct based on partner type
func (pa *PartnerAttributes) GetAttributesForPartnerType(partnerType valueObjects.PartnerType) interface{} {
	switch partnerType {
	case valueObjects.PartnerTypeBank:
		if pa.Bank == nil {
			pa.Bank = &BankPartnerAttributes{}
		}
		return pa.Bank
	case valueObjects.PartnerTypeCooperative:
		if pa.Cooperative == nil {
			pa.Cooperative = &CooperativePartnerAttributes{}
		}
		return pa.Cooperative
	case valueObjects.PartnerTypeFintechSCD:
		if pa.FintechSCD == nil {
			pa.FintechSCD = &FintechSCDAttributes{}
		}
		return pa.FintechSCD
	case valueObjects.PartnerTypeFinanceira:
		if pa.Financeira == nil {
			pa.Financeira = &FinanceiraAttributes{}
		}
		return pa.Financeira
	case valueObjects.PartnerTypePaymentInstitution:
		if pa.PaymentInstitution == nil {
			pa.PaymentInstitution = &PaymentInstitutionAttributes{}
		}
		return pa.PaymentInstitution
	case valueObjects.PartnerTypeSavingsBank:
		if pa.SavingsBank == nil {
			pa.SavingsBank = &SavingsBankAttributes{}
		}
		return pa.SavingsBank
	case valueObjects.PartnerTypeDevelopmentBank:
		if pa.DevelopmentBank == nil {
			pa.DevelopmentBank = &DevelopmentBankAttributes{}
		}
		return pa.DevelopmentBank
	default:
		return nil
	}
}

// Validate checks if the partner attributes are valid for the given partner type
func (pa *PartnerAttributes) Validate(partnerType valueObjects.PartnerType) error {
	switch partnerType {
	case valueObjects.PartnerTypeBank:
		if pa.Bank == nil {
			return errors.New("bank attributes are required for bank partners")
		}
		if pa.Bank.BankCode == "" {
			return errors.New("bank_code is required for bank partners")
		}
	case valueObjects.PartnerTypeCooperative:
		if pa.Cooperative == nil {
			return errors.New("cooperative attributes are required for cooperative partners")
		}
		if pa.Cooperative.CooperativeRegistry == "" {
			return errors.New("cooperative_registry is required for cooperative partners")
		}
	// Add validation for other partner types as needed
	default:
		// For other types, basic validation
		if pa.GetAttributesForPartnerType(partnerType) == nil {
			return errors.New("attributes are required for this partner type")
		}
	}
	return nil
}
