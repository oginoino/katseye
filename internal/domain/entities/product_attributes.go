package entities

import (
	"errors"

	valueObjects "katseye/internal/domain/value_objects"
)

type BaseProductAttributes struct {
	InterestRate          float64                         `json:"interest_rate,omitempty" bson:"interest_rate,omitempty"`
	TermMonths            int                             `json:"term_months,omitempty" bson:"term_months,omitempty"`
	MinAmount             float64                         `json:"min_amount,omitempty" bson:"min_amount,omitempty"`
	MaxAmount             float64                         `json:"max_amount,omitempty" bson:"max_amount,omitempty"`
	ProcessingFee         float64                         `json:"processing_fee,omitempty" bson:"processing_fee,omitempty"`
	EarlyRepaymentAllowed bool                            `json:"early_repayment_allowed,omitempty" bson:"early_repayment_allowed,omitempty"`
	RequiredDocuments     []valueObjects.RequiredDocument `json:"required_documents,omitempty" bson:"required_documents,omitempty"`
	GracePeriodDays       int                             `json:"grace_period_days,omitempty" bson:"grace_period_days,omitempty"`
	IofRate               float64                         `json:"iof_rate,omitempty" bson:"iof_rate,omitempty"`
	CETRate               float64                         `json:"cet_rate,omitempty" bson:"cet_rate,omitempty"`
}

// PersonalLoanAttributes contains personal loan specific attributes
type PersonalLoanAttributes struct {
	BaseProductAttributes  `bson:",inline"`
	CreditAnalysisRequired bool `json:"credit_analysis_required" bson:"credit_analysis_required"`
	SalaryTransferRequired bool `json:"salary_transfer_required,omitempty" bson:"salary_transfer_required,omitempty"`
	InsuranceRequired      bool `json:"insurance_required,omitempty" bson:"insurance_required,omitempty"`
	IncomeProofRequired    bool `json:"income_proof_required" bson:"income_proof_required"`
}

// PayrollLoanAttributes contains payroll loan specific attributes
type PayrollLoanAttributes struct {
	BaseProductAttributes `bson:",inline"`
	DiscountOnPayroll     bool    `json:"discount_on_payroll" bson:"discount_on_payroll"`
	MaximumInstallmentPct float64 `json:"maximum_installment_pct,omitempty" bson:"maximum_installment_pct,omitempty"`
	OnlyForPublicServants bool    `json:"only_for_public_servants,omitempty" bson:"only_for_public_servants,omitempty"`
	RetirementBenefit     bool    `json:"retirement_benefit,omitempty" bson:"retirement_benefit,omitempty"`
}

// CreditCardAttributes contains credit card specific attributes
type CreditCardAttributes struct {
	AnnualFee           float64 `json:"annual_fee" bson:"annual_fee"`
	CreditLimit         float64 `json:"credit_limit,omitempty" bson:"credit_limit,omitempty"`
	GracePeriod         int     `json:"grace_period" bson:"grace_period"`
	RevolvingInterest   float64 `json:"revolving_interest" bson:"revolving_interest"`
	CashAdvanceInterest float64 `json:"cash_advance_interest" bson:"cash_advance_interest"`
	CashAdvanceFee      float64 `json:"cash_advance_fee,omitempty" bson:"cash_advance_fee,omitempty"`
	InternationalUse    bool    `json:"international_use,omitempty" bson:"international_use,omitempty"`
	RewardsProgram      bool    `json:"rewards_program,omitempty" bson:"rewards_program,omitempty"`
	InsuranceIncluded   bool    `json:"insurance_included,omitempty" bson:"insurance_included,omitempty"`
}

// VehicleFinancingAttributes contains vehicle financing specific attributes
type VehicleFinancingAttributes struct {
	BaseProductAttributes `bson:",inline"`
	DownPaymentRequired   bool    `json:"down_payment_required" bson:"down_payment_required"`
	MinDownPaymentPct     float64 `json:"min_down_payment_pct,omitempty" bson:"min_down_payment_pct,omitempty"`
	VehicleAgeLimit       int     `json:"vehicle_age_limit,omitempty" bson:"vehicle_age_limit,omitempty"`
	InsuranceRequired     bool    `json:"insurance_required" bson:"insurance_required"`
	BalloonPayment        bool    `json:"balloon_payment,omitempty" bson:"balloon_payment,omitempty"`
}

// MortgageLoanAttributes contains mortgage loan specific attributes
type MortgageLoanAttributes struct {
	BaseProductAttributes     `bson:",inline"`
	PropertyValue             float64 `json:"property_value" bson:"property_value"`
	MinDownPaymentPct         float64 `json:"min_down_payment_pct" bson:"min_down_payment_pct"`
	LtvRatio                  float64 `json:"ltv_ratio,omitempty" bson:"ltv_ratio,omitempty"`
	PropertyInsuranceRequired bool    `json:"property_insurance_required" bson:"property_insurance_required"`
	AppraisalRequired         bool    `json:"appraisal_required" bson:"appraisal_required"`
	FixedRatePeriod           int     `json:"fixed_rate_period,omitempty" bson:"fixed_rate_period,omitempty"`
}

// WorkingCapitalAttributes contains working capital loan specific attributes
type WorkingCapitalAttributes struct {
	BaseProductAttributes    `bson:",inline"`
	BusinessAgeRequirement   int     `json:"business_age_requirement,omitempty" bson:"business_age_requirement,omitempty"`
	AnnualRevenueRequirement float64 `json:"annual_revenue_requirement,omitempty" bson:"annual_revenue_requirement,omitempty"`
	CollateralRequired       bool    `json:"collateral_required,omitempty" bson:"collateral_required,omitempty"`
	ReceivablesAsCollateral  bool    `json:"receivables_as_collateral,omitempty" bson:"receivables_as_collateral,omitempty"`
}

// StudentLoanAttributes contains student loan specific attributes
type StudentLoanAttributes struct {
	BaseProductAttributes      `bson:",inline"`
	CourseType                 string `json:"course_type" bson:"course_type"`           // undergraduate, graduate, technical
	InstitutionType            string `json:"institution_type" bson:"institution_type"` // public, private
	GracePeriodAfterGraduation int    `json:"grace_period_after_graduation" bson:"grace_period_after_graduation"`
	CoSignerRequired           bool   `json:"co_signer_required,omitempty" bson:"co_signer_required,omitempty"`
}

// GreenLoanAttributes contains green loan specific attributes
type GreenLoanAttributes struct {
	BaseProductAttributes `bson:",inline"`
	EcoFriendlyCategory   string   `json:"eco_friendly_category" bson:"eco_friendly_category"`
	CertificationRequired bool     `json:"certification_required" bson:"certification_required"`
	GovernmentIncentive   bool     `json:"government_incentive,omitempty" bson:"government_incentive,omitempty"`
	EligibleProjects      []string `json:"eligible_projects" bson:"eligible_projects"`
}

// ProductAttributes is a union type that can hold any product-specific attributes
type ProductAttributes struct {
	PersonalLoan     *PersonalLoanAttributes     `json:"personal_loan,omitempty" bson:"personal_loan,omitempty"`
	PayrollLoan      *PayrollLoanAttributes      `json:"payroll_loan,omitempty" bson:"payroll_loan,omitempty"`
	CreditCard       *CreditCardAttributes       `json:"credit_card,omitempty" bson:"credit_card,omitempty"`
	VehicleFinancing *VehicleFinancingAttributes `json:"vehicle_financing,omitempty" bson:"vehicle_financing,omitempty"`
	MortgageLoan     *MortgageLoanAttributes     `json:"mortgage_loan,omitempty" bson:"mortgage_loan,omitempty"`
	WorkingCapital   *WorkingCapitalAttributes   `json:"working_capital,omitempty" bson:"working_capital,omitempty"`
	StudentLoan      *StudentLoanAttributes      `json:"student_loan,omitempty" bson:"student_loan,omitempty"`
	GreenLoan        *GreenLoanAttributes        `json:"green_loan,omitempty" bson:"green_loan,omitempty"`
	OverdraftCredit  *BaseProductAttributes      `json:"overdraft_credit,omitempty" bson:"overdraft_credit,omitempty"`
	SecuredLoan      *BaseProductAttributes      `json:"secured_loan,omitempty" bson:"secured_loan,omitempty"`
	MicrocreditLoan  *BaseProductAttributes      `json:"microcredit_loan,omitempty" bson:"microcredit_loan,omitempty"`
	FGTSLoan         *BaseProductAttributes      `json:"fgts_loan,omitempty" bson:"fgts_loan,omitempty"`
	IRPFLoan         *BaseProductAttributes      `json:"irpf_loan,omitempty" bson:"irpf_loan,omitempty"`
	ReceivablesAdvance      *BaseProductAttributes      `json:"receivables_advance,omitempty" bson:"receivables_advance,omitempty"`
	SecuredOverdraft        *BaseProductAttributes      `json:"secured_overdraft,omitempty" bson:"secured_overdraft,omitempty"`
	InvestmentFinancing     *BaseProductAttributes      `json:"investment_financing,omitempty" bson:"investment_financing,omitempty"`
	BNDESLoan               *BaseProductAttributes      `json:"bndes_loan,omitempty" bson:"bndes_loan,omitempty"`
	AgriculturalCredit      *BaseProductAttributes      `json:"agricultural_credit,omitempty" bson:"agricultural_credit,omitempty"`
	LeasingContract         *BaseProductAttributes      `json:"leasing_contract,omitempty" bson:"leasing_contract,omitempty"`
	SolarEnergyLoan         *BaseProductAttributes      `json:"solar_energy_loan,omitempty" bson:"solar_energy_loan,omitempty"`
	FintechLoan             *BaseProductAttributes      `json:"fintech_loan,omitempty" bson:"fintech_loan,omitempty"`
	MicrocreditSolidaryLoan *BaseProductAttributes      `json:"microcredit_solidary_loan,omitempty" bson:"microcredit_solidary_loan,omitempty"`
	// Add more product types as needed
}

// GetAttributesForProductType returns the appropriate attributes struct based on product type
func (pa *ProductAttributes) GetAttributesForProductType(productType valueObjects.ProductType) interface{} {
	switch productType {
	case valueObjects.ProductTypePersonalLoan:
		if pa.PersonalLoan == nil {
			pa.PersonalLoan = &PersonalLoanAttributes{}
		}
		return pa.PersonalLoan
	case valueObjects.ProductTypePayrollLoan:
		if pa.PayrollLoan == nil {
			pa.PayrollLoan = &PayrollLoanAttributes{}
		}
		return pa.PayrollLoan
	case valueObjects.ProductTypeCreditCard:
		if pa.CreditCard == nil {
			pa.CreditCard = &CreditCardAttributes{}
		}
		return pa.CreditCard
	case valueObjects.ProductTypeVehicleFinancing:
		if pa.VehicleFinancing == nil {
			pa.VehicleFinancing = &VehicleFinancingAttributes{}
		}
		return pa.VehicleFinancing
	case valueObjects.ProductTypeMortgageLoan:
		if pa.MortgageLoan == nil {
			pa.MortgageLoan = &MortgageLoanAttributes{}
		}
		return pa.MortgageLoan
	case valueObjects.ProductTypeWorkingCapitalLoan:
		if pa.WorkingCapital == nil {
			pa.WorkingCapital = &WorkingCapitalAttributes{}
		}
		return pa.WorkingCapital
	case valueObjects.ProductTypeStudentLoan:
		if pa.StudentLoan == nil {
			pa.StudentLoan = &StudentLoanAttributes{}
		}
		return pa.StudentLoan
	case valueObjects.ProductTypeGreenLoan:
		if pa.GreenLoan == nil {
			pa.GreenLoan = &GreenLoanAttributes{}
		}
		return pa.GreenLoan
	case valueObjects.ProductTypeOverdraftCredit:
		if pa.OverdraftCredit == nil {
			pa.OverdraftCredit = &BaseProductAttributes{}
		}
		return pa.OverdraftCredit
	case valueObjects.ProductTypeSecuredLoan:
		if pa.SecuredLoan == nil {
			pa.SecuredLoan = &BaseProductAttributes{}
		}
		return pa.SecuredLoan
	case valueObjects.ProductTypeMicrocreditLoan:
		if pa.MicrocreditLoan == nil {
			pa.MicrocreditLoan = &BaseProductAttributes{}
		}
		return pa.MicrocreditLoan
	case valueObjects.ProductTypeFGTSLoan:
		if pa.FGTSLoan == nil {
			pa.FGTSLoan = &BaseProductAttributes{}
		}
		return pa.FGTSLoan
	case valueObjects.ProductTypeIRPFLoan:
		if pa.IRPFLoan == nil {
			pa.IRPFLoan = &BaseProductAttributes{}
		}
		return pa.IRPFLoan
	case valueObjects.ProductTypeReceivablesAdvance:
		if pa.ReceivablesAdvance == nil {
			pa.ReceivablesAdvance = &BaseProductAttributes{}
		}
		return pa.ReceivablesAdvance
	case valueObjects.ProductTypeSecuredOverdraft:
		if pa.SecuredOverdraft == nil {
			pa.SecuredOverdraft = &BaseProductAttributes{}
		}
		return pa.SecuredOverdraft
	case valueObjects.ProductTypeInvestmentFin:
		if pa.InvestmentFinancing == nil {
			pa.InvestmentFinancing = &BaseProductAttributes{}
		}
		return pa.InvestmentFinancing
	case valueObjects.ProductTypeBNDESLoan:
		if pa.BNDESLoan == nil {
			pa.BNDESLoan = &BaseProductAttributes{}
		}
		return pa.BNDESLoan
	case valueObjects.ProductTypeAgriculturalCredit:
		if pa.AgriculturalCredit == nil {
			pa.AgriculturalCredit = &BaseProductAttributes{}
		}
		return pa.AgriculturalCredit
	case valueObjects.ProductTypeLeasingContract:
		if pa.LeasingContract == nil {
			pa.LeasingContract = &BaseProductAttributes{}
		}
		return pa.LeasingContract
	case valueObjects.ProductTypeSolarEnergyLoan:
		if pa.SolarEnergyLoan == nil {
			pa.SolarEnergyLoan = &BaseProductAttributes{}
		}
		return pa.SolarEnergyLoan
	case valueObjects.ProductTypeFintechLoan:
		if pa.FintechLoan == nil {
			pa.FintechLoan = &BaseProductAttributes{}
		}
		return pa.FintechLoan
	case valueObjects.ProductTypeMicrocreditSolidaryLoan:
		if pa.MicrocreditSolidaryLoan == nil {
			pa.MicrocreditSolidaryLoan = &BaseProductAttributes{}
		}
		return pa.MicrocreditSolidaryLoan
	default:
		return nil
	}
}

// Validate checks if the product attributes are valid for the given product type
func (pa *ProductAttributes) Validate(productType valueObjects.ProductType) error {
	switch productType {
	case valueObjects.ProductTypePersonalLoan:
		if pa.PersonalLoan == nil {
			return errors.New("personal_loan attributes are required for personal loan products")
		}
		if pa.PersonalLoan.InterestRate <= 0 {
			return errors.New("interest_rate is required for personal loan products")
		}
		// Validate required documents
		if err := valueObjects.ValidateDocumentSet(pa.PersonalLoan.RequiredDocuments); err != nil {
			return errors.New("invalid required documents for personal loan: " + err.Error())
		}

	case valueObjects.ProductTypePayrollLoan:
		if pa.PayrollLoan == nil {
			return errors.New("payroll_loan attributes are required for payroll loan products")
		}
		if err := valueObjects.ValidateDocumentSet(pa.PayrollLoan.RequiredDocuments); err != nil {
			return errors.New("invalid required documents for payroll loan: " + err.Error())
		}

	case valueObjects.ProductTypeCreditCard:
		if pa.CreditCard == nil {
			return errors.New("credit_card attributes are required for credit card products")
		}
		if pa.CreditCard.AnnualFee < 0 {
			return errors.New("annual_fee is required for credit card products")
		}

	case valueObjects.ProductTypeVehicleFinancing:
		if pa.VehicleFinancing == nil {
			return errors.New("vehicle_financing attributes are required for vehicle financing products")
		}
		if err := valueObjects.ValidateDocumentSet(pa.VehicleFinancing.RequiredDocuments); err != nil {
			return errors.New("invalid required documents for vehicle financing: " + err.Error())
		}

	case valueObjects.ProductTypeMortgageLoan:
		if pa.MortgageLoan == nil {
			return errors.New("mortgage_loan attributes are required for mortgage loan products")
		}
		if pa.MortgageLoan.PropertyValue <= 0 {
			return errors.New("property_value is required for mortgage loan products")
		}
		if err := valueObjects.ValidateDocumentSet(pa.MortgageLoan.RequiredDocuments); err != nil {
			return errors.New("invalid required documents for mortgage loan: " + err.Error())
		}

	// Add validation for other product types as needed
	default:
		// For other types, basic validation
		if pa.GetAttributesForProductType(productType) == nil {
			return errors.New("attributes are required for this product type")
		}
	}
	return nil
}
