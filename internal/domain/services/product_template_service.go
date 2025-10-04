package services

import (
	"sort"

	"katseye/internal/domain/entities"
	valueobjects "katseye/internal/domain/value_objects"
)

type ProductTemplate struct {
	Type        valueobjects.ProductType
	DisplayName string
	Category    string
	Attributes  entities.ProductAttributes
}

type ProductTemplateService struct {
	templates map[valueobjects.ProductType]ProductTemplate
	ordered   []ProductTemplate
}

func NewProductTemplateService() *ProductTemplateService {
	templates := []ProductTemplate{
		buildPersonalLoanTemplate(),
		buildPayrollLoanTemplate(),
		buildCreditCardTemplate(),
		buildOverdraftCreditTemplate(),
		buildVehicleFinancingTemplate(),
		buildMortgageLoanTemplate(),
		buildSecuredLoanTemplate(),
		buildMicrocreditLoanTemplate(),
		buildFGTSLoanTemplate(),
		buildIRPFLoanTemplate(),
		buildWorkingCapitalTemplate(),
		buildReceivablesAdvanceTemplate(),
		buildSecuredOverdraftTemplate(),
		buildInvestmentFinancingTemplate(),
		buildBNDESLoanTemplate(),
		buildAgriculturalCreditTemplate(),
		buildLeasingContractTemplate(),
		buildStudentLoanTemplate(),
		buildGreenLoanTemplate(),
		buildSolarEnergyLoanTemplate(),
		buildFintechLoanTemplate(),
		buildMicrocreditSolidaryLoanTemplate(),
	}

	ordered := make([]ProductTemplate, len(templates))
	copy(ordered, templates)
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].DisplayName < ordered[j].DisplayName
	})

	templateMap := make(map[valueobjects.ProductType]ProductTemplate, len(templates))
	for _, template := range templates {
		templateMap[template.Type] = template
	}

	return &ProductTemplateService{
		templates: templateMap,
		ordered:   ordered,
	}
}

func (s *ProductTemplateService) ListTemplates() []ProductTemplate {
	if s == nil {
		return nil
	}
	result := make([]ProductTemplate, len(s.ordered))
	copy(result, s.ordered)
	return result
}

func (s *ProductTemplateService) GetTemplate(productType valueobjects.ProductType) (ProductTemplate, bool) {
	if s == nil {
		return ProductTemplate{}, false
	}
	template, ok := s.templates[productType]
	return template, ok
}

func buildPersonalLoanTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypePersonalLoan,
		DisplayName: "Empréstimo pessoal",
		Category:    "credit",
		Attributes: entities.ProductAttributes{
			PersonalLoan: &entities.PersonalLoanAttributes{
				BaseProductAttributes: entities.BaseProductAttributes{
					InterestRate:          2.3,
					TermMonths:            36,
					MinAmount:             1000,
					MaxAmount:             50000,
					ProcessingFee:         1.5,
					EarlyRepaymentAllowed: true,
					RequiredDocuments: []valueobjects.RequiredDocument{
						valueobjects.DocumentIDProof,
						valueobjects.DocumentIncomeProof,
					},
					GracePeriodDays: 30,
					IofRate:         0.38,
					CETRate:         2.9,
				},
				CreditAnalysisRequired: true,
				SalaryTransferRequired: false,
				InsuranceRequired:      false,
				IncomeProofRequired:    true,
			},
		},
	}
}

func buildPayrollLoanTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypePayrollLoan,
		DisplayName: "Crédito consignado",
		Category:    "credit",
		Attributes: entities.ProductAttributes{
			PayrollLoan: &entities.PayrollLoanAttributes{
				BaseProductAttributes: entities.BaseProductAttributes{
					InterestRate:          1.9,
					TermMonths:            60,
					MinAmount:             500,
					MaxAmount:             80000,
					ProcessingFee:         0,
					EarlyRepaymentAllowed: true,
					RequiredDocuments: []valueobjects.RequiredDocument{
						valueobjects.DocumentIDProof,
						valueobjects.DocumentPayStub,
					},
				},
				DiscountOnPayroll:     true,
				MaximumInstallmentPct: 35,
				OnlyForPublicServants: false,
				RetirementBenefit:     true,
			},
		},
	}
}

func buildCreditCardTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeCreditCard,
		DisplayName: "Cartão de crédito",
		Category:    "credit",
		Attributes: entities.ProductAttributes{
			CreditCard: &entities.CreditCardAttributes{
				AnnualFee:           150,
				CreditLimit:         10000,
				GracePeriod:         40,
				RevolvingInterest:   12.5,
				CashAdvanceInterest: 9.9,
				CashAdvanceFee:      3.9,
				InternationalUse:    true,
				RewardsProgram:      true,
				InsuranceIncluded:   false,
			},
		},
	}
}

func buildOverdraftCreditTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeOverdraftCredit,
		DisplayName: "Cheque especial",
		Category:    "credit",
		Attributes: entities.ProductAttributes{
			OverdraftCredit: &entities.BaseProductAttributes{
				InterestRate:          3.9,
				TermMonths:            12,
				MinAmount:             100,
				MaxAmount:             20000,
				ProcessingFee:         0,
				EarlyRepaymentAllowed: true,
				RequiredDocuments: []valueobjects.RequiredDocument{
					valueobjects.DocumentIDProof,
					valueobjects.DocumentBankStatement,
				},
			},
		},
	}
}

func buildVehicleFinancingTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeVehicleFinancing,
		DisplayName: "Financiamento de veículos",
		Category:    "credit",
		Attributes: entities.ProductAttributes{
			VehicleFinancing: &entities.VehicleFinancingAttributes{
				BaseProductAttributes: entities.BaseProductAttributes{
					InterestRate:          1.6,
					TermMonths:            48,
					MinAmount:             5000,
					MaxAmount:             150000,
					ProcessingFee:         1.2,
					EarlyRepaymentAllowed: true,
					RequiredDocuments: []valueobjects.RequiredDocument{
						valueobjects.DocumentIDProof,
						valueobjects.DocumentVehicleInvoice,
					},
				},
				DownPaymentRequired: true,
				MinDownPaymentPct:   20,
				VehicleAgeLimit:     8,
				InsuranceRequired:   true,
				BalloonPayment:      false,
			},
		},
	}
}

func buildMortgageLoanTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeMortgageLoan,
		DisplayName: "Crédito imobiliário",
		Category:    "credit",
		Attributes: entities.ProductAttributes{
			MortgageLoan: &entities.MortgageLoanAttributes{
				BaseProductAttributes: entities.BaseProductAttributes{
					InterestRate:          1.4,
					TermMonths:            240,
					MinAmount:             50000,
					MaxAmount:             1000000,
					ProcessingFee:         0.9,
					EarlyRepaymentAllowed: true,
					RequiredDocuments: []valueobjects.RequiredDocument{
						valueobjects.DocumentIDProof,
						valueobjects.DocumentPropertyDeed,
						valueobjects.DocumentFinancialStatements,
					},
				},
				PropertyValue:             300000,
				MinDownPaymentPct:         20,
				LtvRatio:                  80,
				PropertyInsuranceRequired: true,
				AppraisalRequired:         true,
				FixedRatePeriod:           36,
			},
		},
	}
}

func buildSecuredLoanTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeSecuredLoan,
		DisplayName: "Crédito com garantia",
		Category:    "credit",
		Attributes: entities.ProductAttributes{
			SecuredLoan: &entities.BaseProductAttributes{
				InterestRate:          1.7,
				TermMonths:            60,
				MinAmount:             10000,
				MaxAmount:             300000,
				ProcessingFee:         0.5,
				EarlyRepaymentAllowed: true,
				RequiredDocuments: []valueobjects.RequiredDocument{
					valueobjects.DocumentIDProof,
					valueobjects.DocumentPropertyDeed,
				},
			},
		},
	}
}

func buildMicrocreditLoanTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeMicrocreditLoan,
		DisplayName: "Microcrédito",
		Category:    "credit",
		Attributes: entities.ProductAttributes{
			MicrocreditLoan: &entities.BaseProductAttributes{
				InterestRate:          2.1,
				TermMonths:            24,
				MinAmount:             500,
				MaxAmount:             15000,
				ProcessingFee:         0.3,
				EarlyRepaymentAllowed: true,
				RequiredDocuments: []valueobjects.RequiredDocument{
					valueobjects.DocumentIDProof,
					valueobjects.DocumentIncomeProof,
				},
			},
		},
	}
}

func buildFGTSLoanTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeFGTSLoan,
		DisplayName: "Antecipação FGTS",
		Category:    "credit",
		Attributes: entities.ProductAttributes{
			FGTSLoan: &entities.BaseProductAttributes{
				InterestRate:          1.2,
				TermMonths:            24,
				MinAmount:             500,
				MaxAmount:             30000,
				ProcessingFee:         0,
				EarlyRepaymentAllowed: true,
				RequiredDocuments: []valueobjects.RequiredDocument{
					valueobjects.DocumentIDProof,
					valueobjects.DocumentEmploymentLetter,
				},
			},
		},
	}
}

func buildIRPFLoanTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeIRPFLoan,
		DisplayName: "Antecipação IRPF",
		Category:    "credit",
		Attributes: entities.ProductAttributes{
			IRPFLoan: &entities.BaseProductAttributes{
				InterestRate:          1.4,
				TermMonths:            12,
				MinAmount:             500,
				MaxAmount:             20000,
				ProcessingFee:         0,
				EarlyRepaymentAllowed: true,
				RequiredDocuments: []valueobjects.RequiredDocument{
					valueobjects.DocumentIDProof,
					valueobjects.DocumentTaxReturn,
				},
			},
		},
	}
}

func buildWorkingCapitalTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeWorkingCapitalLoan,
		DisplayName: "Capital de giro",
		Category:    "business",
		Attributes: entities.ProductAttributes{
			WorkingCapital: &entities.WorkingCapitalAttributes{
				BaseProductAttributes: entities.BaseProductAttributes{
					InterestRate:          1.8,
					TermMonths:            24,
					MinAmount:             10000,
					MaxAmount:             500000,
					ProcessingFee:         0.8,
					EarlyRepaymentAllowed: true,
					RequiredDocuments: []valueobjects.RequiredDocument{
						valueobjects.DocumentCNPJ,
						valueobjects.DocumentFinancialStatements,
					},
				},
				BusinessAgeRequirement:   12,
				AnnualRevenueRequirement: 200000,
				CollateralRequired:       false,
				ReceivablesAsCollateral:  false,
			},
		},
	}
}

func buildReceivablesAdvanceTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeReceivablesAdvance,
		DisplayName: "Antecipação de recebíveis",
		Category:    "business",
		Attributes: entities.ProductAttributes{
			ReceivablesAdvance: &entities.BaseProductAttributes{
				InterestRate:          1.3,
				TermMonths:            6,
				MinAmount:             2000,
				MaxAmount:             250000,
				ProcessingFee:         0.4,
				EarlyRepaymentAllowed: true,
				RequiredDocuments: []valueobjects.RequiredDocument{
					valueobjects.DocumentCNPJ,
					valueobjects.DocumentFinancialStatements,
					valueobjects.DocumentBankStatement,
				},
			},
		},
	}
}

func buildSecuredOverdraftTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeSecuredOverdraft,
		DisplayName: "Cheque especial garantido",
		Category:    "business",
		Attributes: entities.ProductAttributes{
			SecuredOverdraft: &entities.BaseProductAttributes{
				InterestRate:          2.4,
				TermMonths:            18,
				MinAmount:             1000,
				MaxAmount:             100000,
				ProcessingFee:         0.2,
				EarlyRepaymentAllowed: true,
				RequiredDocuments: []valueobjects.RequiredDocument{
					valueobjects.DocumentIDProof,
					valueobjects.DocumentBankStatement,
					valueobjects.DocumentPropertyDeed,
				},
			},
		},
	}
}

func buildInvestmentFinancingTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeInvestmentFin,
		DisplayName: "Financiamento de investimentos",
		Category:    "business",
		Attributes: entities.ProductAttributes{
			InvestmentFinancing: &entities.BaseProductAttributes{
				InterestRate:          1.6,
				TermMonths:            84,
				MinAmount:             50000,
				MaxAmount:             800000,
				ProcessingFee:         0.6,
				EarlyRepaymentAllowed: true,
				RequiredDocuments: []valueobjects.RequiredDocument{
					valueobjects.DocumentCNPJ,
					valueobjects.DocumentBusinessLicense,
					valueobjects.DocumentFinancialStatements,
				},
			},
		},
	}
}

func buildBNDESLoanTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeBNDESLoan,
		DisplayName: "Linha BNDES",
		Category:    "business",
		Attributes: entities.ProductAttributes{
			BNDESLoan: &entities.BaseProductAttributes{
				InterestRate:          1.1,
				TermMonths:            120,
				MinAmount:             100000,
				MaxAmount:             2000000,
				ProcessingFee:         0.5,
				EarlyRepaymentAllowed: true,
				RequiredDocuments: []valueobjects.RequiredDocument{
					valueobjects.DocumentCNPJ,
					valueobjects.DocumentFinancialStatements,
					valueobjects.DocumentBusinessLicense,
				},
			},
		},
	}
}

func buildAgriculturalCreditTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeAgriculturalCredit,
		DisplayName: "Crédito agrícola",
		Category:    "agribusiness",
		Attributes: entities.ProductAttributes{
			AgriculturalCredit: &entities.BaseProductAttributes{
				InterestRate:          1.3,
				TermMonths:            48,
				MinAmount:             20000,
				MaxAmount:             700000,
				ProcessingFee:         0.3,
				EarlyRepaymentAllowed: true,
				RequiredDocuments: []valueobjects.RequiredDocument{
					valueobjects.DocumentCNPJ,
					valueobjects.DocumentPropertyDeed,
					valueobjects.DocumentFinancialStatements,
				},
			},
		},
	}
}

func buildLeasingContractTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeLeasingContract,
		DisplayName: "Leasing",
		Category:    "business",
		Attributes: entities.ProductAttributes{
			LeasingContract: &entities.BaseProductAttributes{
				InterestRate:          1.5,
				TermMonths:            60,
				MinAmount:             10000,
				MaxAmount:             500000,
				ProcessingFee:         0.7,
				EarlyRepaymentAllowed: true,
				RequiredDocuments: []valueobjects.RequiredDocument{
					valueobjects.DocumentIDProof,
					valueobjects.DocumentVehicleInvoice,
				},
			},
		},
	}
}

func buildStudentLoanTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeStudentLoan,
		DisplayName: "Crédito estudantil",
		Category:    "education",
		Attributes: entities.ProductAttributes{
			StudentLoan: &entities.StudentLoanAttributes{
				BaseProductAttributes: entities.BaseProductAttributes{
					InterestRate:          0.9,
					TermMonths:            72,
					MinAmount:             1000,
					MaxAmount:             100000,
					ProcessingFee:         0,
					EarlyRepaymentAllowed: false,
					RequiredDocuments: []valueobjects.RequiredDocument{
						valueobjects.DocumentStudentID,
						valueobjects.DocumentEnrollmentProof,
					},
					GracePeriodDays: 0,
				},
				CourseType:                 "graduate",
				InstitutionType:            "private",
				GracePeriodAfterGraduation: 18,
				CoSignerRequired:           false,
			},
		},
	}
}

func buildGreenLoanTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeGreenLoan,
		DisplayName: "Crédito verde",
		Category:    "sustainability",
		Attributes: entities.ProductAttributes{
			GreenLoan: &entities.GreenLoanAttributes{
				BaseProductAttributes: entities.BaseProductAttributes{
					InterestRate:          1.2,
					TermMonths:            60,
					MinAmount:             500,
					MaxAmount:             50000,
					ProcessingFee:         0,
					EarlyRepaymentAllowed: true,
					RequiredDocuments: []valueobjects.RequiredDocument{
						valueobjects.DocumentIDProof,
						valueobjects.DocumentIncomeProof,
					},
					GracePeriodDays: 30,
					IofRate:         0.38,
					CETRate:         2.9,
				},
				EcoFriendlyCategory:   valueobjects.EcoFriendlyCategoryGreen,
				CertificationRequired: valueobjects.CertificationRequiredTrue,
				GovernmentIncentive:   valueobjects.GovernmentIncentiveTrue,
				EligibleProjects:      valueobjects.EligibleProjectsTrue,
			},
		},
	}
}

func buildSolarEnergyLoanTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeSolarEnergyLoan,
		DisplayName: "Crédito energia solar",
		Category:    "sustainability",
		Attributes: entities.ProductAttributes{
			SolarEnergyLoan: &entities.BaseProductAttributes{
				InterestRate:          1.1,
				TermMonths:            72,
				MinAmount:             5000,
				MaxAmount:             200000,
				ProcessingFee:         0.4,
				EarlyRepaymentAllowed: true,
				RequiredDocuments: []valueobjects.RequiredDocument{
					valueobjects.DocumentIDProof,
					valueobjects.DocumentPropertyDeed,
					valueobjects.DocumentHomeInsurance,
				},
			},
		},
	}
}

func buildFintechLoanTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeFintechLoan,
		DisplayName: "Crédito fintech",
		Category:    "innovation",
		Attributes: entities.ProductAttributes{
			FintechLoan: &entities.BaseProductAttributes{
				InterestRate:          1.7,
				TermMonths:            36,
				MinAmount:             1000,
				MaxAmount:             100000,
				ProcessingFee:         0.2,
				EarlyRepaymentAllowed: true,
				RequiredDocuments: []valueobjects.RequiredDocument{
					valueobjects.DocumentIDProof,
					valueobjects.DocumentBankStatement,
				},
			},
		},
	}
}

func buildMicrocreditSolidaryLoanTemplate() ProductTemplate {
	return ProductTemplate{
		Type:        valueobjects.ProductTypeMicrocreditSolidaryLoan,
		DisplayName: "Microcrédito solidário",
		Category:    "social",
		Attributes: entities.ProductAttributes{
			MicrocreditSolidaryLoan: &entities.BaseProductAttributes{
				InterestRate:          1.5,
				TermMonths:            18,
				MinAmount:             300,
				MaxAmount:             5000,
				ProcessingFee:         0,
				EarlyRepaymentAllowed: true,
				RequiredDocuments: []valueobjects.RequiredDocument{
					valueobjects.DocumentIDProof,
					valueobjects.DocumentIncomeProof,
				},
			},
		},
	}
}
