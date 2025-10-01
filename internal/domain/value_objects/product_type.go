package valueobjects

import "errors"

type ProductType string

const (
	// Personal Finance Products
	ProductTypePersonalLoan     ProductType = "personal_loan"
	ProductTypePayrollLoan      ProductType = "payroll_loan"
	ProductTypeCreditCard       ProductType = "credit_card"
	ProductTypeOverdraftCredit  ProductType = "overdraft_credit"
	ProductTypeVehicleFinancing ProductType = "vehicle_financing"
	ProductTypeMortgageLoan     ProductType = "mortgage_loan"
	ProductTypeSecuredLoan      ProductType = "secured_loan"
	ProductTypeMicrocreditLoan  ProductType = "microcredit_loan"
	ProductTypeFGTSLoan         ProductType = "fgts_loan"
	ProductTypeIRPFLoan         ProductType = "irpf_loan"

	// Business Products
	ProductTypeWorkingCapitalLoan ProductType = "working_capital_loan"
	ProductTypeReceivablesAdvance ProductType = "receivables_advance"
	ProductTypeSecuredOverdraft   ProductType = "secured_overdraft"
	ProductTypeInvestmentFin      ProductType = "investment_financing"
	ProductTypeBNDESLoan          ProductType = "bndes_loan"
	ProductTypeAgriculturalCredit ProductType = "agricultural_credit"
	ProductTypeLeasingContract    ProductType = "leasing_contract"

	// Specialized Products
	ProductTypeStudentLoan             ProductType = "student_loan"
	ProductTypeGreenLoan               ProductType = "green_loan"
	ProductTypeSolarEnergyLoan         ProductType = "solar_energy_loan"
	ProductTypeFintechLoan             ProductType = "fintech_loan"
	ProductTypeMicrocreditSolidaryLoan ProductType = "microcredit_solidary_loan"
)

var validProductTypes = map[ProductType]bool{
	// Personal Finance Products
	ProductTypePersonalLoan:     true,
	ProductTypePayrollLoan:      true,
	ProductTypeCreditCard:       true,
	ProductTypeOverdraftCredit:  true,
	ProductTypeVehicleFinancing: true,
	ProductTypeMortgageLoan:     true,
	ProductTypeSecuredLoan:      true,
	ProductTypeMicrocreditLoan:  true,
	ProductTypeFGTSLoan:         true,
	ProductTypeIRPFLoan:         true,

	// Business Products
	ProductTypeWorkingCapitalLoan: true,
	ProductTypeReceivablesAdvance: true,
	ProductTypeSecuredOverdraft:   true,
	ProductTypeInvestmentFin:      true,
	ProductTypeBNDESLoan:          true,
	ProductTypeAgriculturalCredit: true,
	ProductTypeLeasingContract:    true,

	// Specialized Products
	ProductTypeStudentLoan:             true,
	ProductTypeGreenLoan:               true,
	ProductTypeSolarEnergyLoan:         true,
	ProductTypeFintechLoan:             true,
	ProductTypeMicrocreditSolidaryLoan: true,
}

var ErrInvalidProductType = errors.New("invalid product type")

// NewProductType creates a validated ProductType
func NewProductType(value string) (ProductType, error) {
	pt := ProductType(value)
	if err := pt.Validate(); err != nil {
		return "", err
	}
	return pt, nil
}

// Validate checks if the ProductType is valid
func (pt ProductType) Validate() error {
	if _, exists := validProductTypes[pt]; !exists {
		return ErrInvalidProductType
	}
	return nil
}

// String returns the string representation
func (pt ProductType) String() string {
	return string(pt)
}

// IsPersonalFinance checks if the product type is a personal finance product
func (pt ProductType) IsPersonalFinance() bool {
	personalFinanceTypes := []ProductType{
		ProductTypePersonalLoan,
		ProductTypePayrollLoan,
		ProductTypeCreditCard,
		ProductTypeOverdraftCredit,
		ProductTypeVehicleFinancing,
		ProductTypeMortgageLoan,
		ProductTypeSecuredLoan,
		ProductTypeMicrocreditLoan,
		ProductTypeFGTSLoan,
		ProductTypeIRPFLoan,
	}
	for _, pft := range personalFinanceTypes {
		if pt == pft {
			return true
		}
	}
	return false
}

// IsBusinessProduct checks if the product type is a business product
func (pt ProductType) IsBusinessProduct() bool {
	businessProductTypes := []ProductType{
		ProductTypeWorkingCapitalLoan,
		ProductTypeReceivablesAdvance,
		ProductTypeSecuredOverdraft,
		ProductTypeInvestmentFin,
		ProductTypeBNDESLoan,
		ProductTypeAgriculturalCredit,
		ProductTypeLeasingContract,
	}
	for _, bpt := range businessProductTypes {
		if pt == bpt {
			return true
		}
	}
	return false
}

// IsSpecializedProduct checks if the product type is a specialized product
func (pt ProductType) IsSpecializedProduct() bool {
	specializedProductTypes := []ProductType{
		ProductTypeStudentLoan,
		ProductTypeGreenLoan,
		ProductTypeSolarEnergyLoan,
		ProductTypeFintechLoan,
		ProductTypeMicrocreditSolidaryLoan,
	}
	for _, spt := range specializedProductTypes {
		if pt == spt {
			return true
		}
	}
	return false
}

// GetRiskLevel returns the risk level associated with the product type
func (pt ProductType) GetRiskLevel() RiskLevel {
	switch pt {
	// Low Risk Products
	case ProductTypeLeasingContract:
		return RiskLevelLow

	// Medium Risk Products
	case ProductTypePersonalLoan, ProductTypePayrollLoan, ProductTypeCreditCard,
		ProductTypeOverdraftCredit, ProductTypeVehicleFinancing, ProductTypeMortgageLoan,
		ProductTypeSecuredLoan, ProductTypeMicrocreditLoan, ProductTypeFGTSLoan,
		ProductTypeIRPFLoan, ProductTypeWorkingCapitalLoan, ProductTypeReceivablesAdvance,
		ProductTypeSecuredOverdraft, ProductTypeInvestmentFin, ProductTypeBNDESLoan,
		ProductTypeAgriculturalCredit:
		return RiskLevelMedium
	// High Risk Products
	case ProductTypeStudentLoan, ProductTypeGreenLoan, ProductTypeSolarEnergyLoan,
		ProductTypeFintechLoan, ProductTypeMicrocreditSolidaryLoan:
		return RiskLevelHigh
	default:
		return RiskLevelMedium // Default to medium risk if unknown
	}
}
