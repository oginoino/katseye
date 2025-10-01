package valueobjects

import (
	"errors"
	"strings"
)

type RequiredDocument string

const (
	// Identity Documents
	DocumentIDProof       RequiredDocument = "id_proof"
	DocumentCPF           RequiredDocument = "cpf"
	DocumentCNPJ          RequiredDocument = "cnpj"
	DocumentRG            RequiredDocument = "rg"
	DocumentPassport      RequiredDocument = "passport"
	DocumentDriverLicense RequiredDocument = "driver_license"

	// Income Documents
	DocumentIncomeProof      RequiredDocument = "income_proof"
	DocumentPayStub          RequiredDocument = "pay_stub"
	DocumentTaxReturn        RequiredDocument = "tax_return"
	DocumentBankStatement    RequiredDocument = "bank_statement"
	DocumentEmploymentLetter RequiredDocument = "employment_letter"
	DocumentProfitLoss       RequiredDocument = "profit_loss"

	// Address Documents
	DocumentAddressProof    RequiredDocument = "address_proof"
	DocumentUtilityBill     RequiredDocument = "utility_bill"
	DocumentRentalAgreement RequiredDocument = "rental_agreement"
	DocumentPropertyDeed    RequiredDocument = "property_deed"

	// Business Documents
	DocumentBusinessLicense         RequiredDocument = "business_license"
	DocumentArticlesOfIncorporation RequiredDocument = "articles_of_incorporation"
	DocumentCommercialLicense       RequiredDocument = "commercial_license"
	DocumentFinancialStatements     RequiredDocument = "financial_statements"

	// Property Documents
	DocumentPropertyTax     RequiredDocument = "property_tax"
	DocumentHomeInsurance   RequiredDocument = "home_insurance"
	DocumentAppraisalReport RequiredDocument = "appraisal_report"

	// Vehicle Documents
	DocumentVehicleRegistration RequiredDocument = "vehicle_registration"
	DocumentVehicleInsurance    RequiredDocument = "vehicle_insurance"
	DocumentVehicleInvoice      RequiredDocument = "vehicle_invoice"

	// Educational Documents
	DocumentStudentID       RequiredDocument = "student_id"
	DocumentEnrollmentProof RequiredDocument = "enrollment_proof"
	DocumentAcademicRecords RequiredDocument = "academic_records"

	// Special Documents
	DocumentMarriageCertificate RequiredDocument = "marriage_certificate"
	DocumentBirthCertificate    RequiredDocument = "birth_certificate"
	DocumentMilitaryID          RequiredDocument = "military_id"
	DocumentProfessionalLicense RequiredDocument = "professional_license"
)

var (
	ErrInvalidRequiredDocument = errors.New("invalid required document")

	validRequiredDocuments = map[RequiredDocument]bool{
		// Identity Documents
		DocumentIDProof:       true,
		DocumentCPF:           true,
		DocumentCNPJ:          true,
		DocumentRG:            true,
		DocumentPassport:      true,
		DocumentDriverLicense: true,

		// Income Documents
		DocumentIncomeProof:      true,
		DocumentPayStub:          true,
		DocumentTaxReturn:        true,
		DocumentBankStatement:    true,
		DocumentEmploymentLetter: true,
		DocumentProfitLoss:       true,

		// Address Documents
		DocumentAddressProof:    true,
		DocumentUtilityBill:     true,
		DocumentRentalAgreement: true,
		DocumentPropertyDeed:    true,

		// Business Documents
		DocumentBusinessLicense:         true,
		DocumentArticlesOfIncorporation: true,
		DocumentCommercialLicense:       true,
		DocumentFinancialStatements:     true,

		// Property Documents
		DocumentPropertyTax:     true,
		DocumentHomeInsurance:   true,
		DocumentAppraisalReport: true,

		// Vehicle Documents
		DocumentVehicleRegistration: true,
		DocumentVehicleInsurance:    true,
		DocumentVehicleInvoice:      true,

		// Educational Documents
		DocumentStudentID:       true,
		DocumentEnrollmentProof: true,
		DocumentAcademicRecords: true,

		// Special Documents
		DocumentMarriageCertificate: true,
		DocumentBirthCertificate:    true,
		DocumentMilitaryID:          true,
		DocumentProfessionalLicense: true,
	}

	// DocumentCategories groups documents by type for easier validation
	DocumentCategories = map[string][]RequiredDocument{
		"identity": {
			DocumentIDProof,
			DocumentCPF,
			DocumentCNPJ,
			DocumentRG,
			DocumentPassport,
			DocumentDriverLicense,
		},
		"income": {
			DocumentIncomeProof,
			DocumentPayStub,
			DocumentTaxReturn,
			DocumentBankStatement,
			DocumentEmploymentLetter,
			DocumentProfitLoss,
		},
		"address": {
			DocumentAddressProof,
			DocumentUtilityBill,
			DocumentRentalAgreement,
			DocumentPropertyDeed,
		},
		"business": {
			DocumentBusinessLicense,
			DocumentArticlesOfIncorporation,
			DocumentCommercialLicense,
			DocumentFinancialStatements,
		},
		"property": {
			DocumentPropertyTax,
			DocumentHomeInsurance,
			DocumentAppraisalReport,
		},
		"vehicle": {
			DocumentVehicleRegistration,
			DocumentVehicleInsurance,
			DocumentVehicleInvoice,
		},
		"education": {
			DocumentStudentID,
			DocumentEnrollmentProof,
			DocumentAcademicRecords,
		},
		"special": {
			DocumentMarriageCertificate,
			DocumentBirthCertificate,
			DocumentMilitaryID,
			DocumentProfessionalLicense,
		},
	}
)

// NewRequiredDocument creates a validated RequiredDocument
func NewRequiredDocument(value string) (RequiredDocument, error) {
	doc := RequiredDocument(strings.ToLower(value))
	if err := doc.Validate(); err != nil {
		return "", err
	}
	return doc, nil
}

// Validate checks if the required document is valid
func (rd RequiredDocument) Validate() error {
	if !validRequiredDocuments[rd] {
		return ErrInvalidRequiredDocument
	}
	return nil
}

// String returns the string representation
func (rd RequiredDocument) String() string {
	return string(rd)
}

// IsIdentityDocument checks if the document is an identity document
func (rd RequiredDocument) IsIdentityDocument() bool {
	for _, doc := range DocumentCategories["identity"] {
		if rd == doc {
			return true
		}
	}
	return false
}

// IsIncomeDocument checks if the document is an income document
func (rd RequiredDocument) IsIncomeDocument() bool {
	for _, doc := range DocumentCategories["income"] {
		if rd == doc {
			return true
		}
	}
	return false
}

// IsAddressDocument checks if the document is an address document
func (rd RequiredDocument) IsAddressDocument() bool {
	for _, doc := range DocumentCategories["address"] {
		if rd == doc {
			return true
		}
	}
	return false
}

// GetDocumentCategory returns the category of the document
func (rd RequiredDocument) GetDocumentCategory() string {
	for category, documents := range DocumentCategories {
		for _, doc := range documents {
			if rd == doc {
				return category
			}
		}
	}
	return "unknown"
}

// ValidateDocumentSet validates a set of required documents
func ValidateDocumentSet(documents []RequiredDocument) error {
	for _, doc := range documents {
		if err := doc.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// GetDocumentsByCategory returns all documents in a specific category
func GetDocumentsByCategory(category string) []RequiredDocument {
	return DocumentCategories[category]
}

// CommonDocumentSets provides predefined document sets for common use cases
var CommonDocumentSets = map[string][]RequiredDocument{
	"personal_loan_basic": {
		DocumentIDProof,
		DocumentCPF,
		DocumentIncomeProof,
		DocumentAddressProof,
	},
	"personal_loan_complete": {
		DocumentIDProof,
		DocumentCPF,
		DocumentIncomeProof,
		DocumentPayStub,
		DocumentBankStatement,
		DocumentAddressProof,
	},
	"business_loan": {
		DocumentCNPJ,
		DocumentBusinessLicense,
		DocumentFinancialStatements,
		DocumentBankStatement,
		DocumentAddressProof,
	},
	"mortgage_loan": {
		DocumentIDProof,
		DocumentCPF,
		DocumentIncomeProof,
		DocumentAddressProof,
		DocumentPropertyDeed,
		DocumentPropertyTax,
		DocumentAppraisalReport,
	},
	"vehicle_financing": {
		DocumentIDProof,
		DocumentCPF,
		DocumentIncomeProof,
		DocumentAddressProof,
		DocumentDriverLicense,
		DocumentVehicleRegistration,
	},
	"student_loan": {
		DocumentIDProof,
		DocumentCPF,
		DocumentStudentID,
		DocumentEnrollmentProof,
		DocumentAddressProof,
	},
}
