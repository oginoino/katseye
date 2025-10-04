package entities

import (
	"errors"
	"fmt"
	"strings"
	"time"

	valueobjects "katseye/internal/domain/value_objects"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrConsumerNil                    = errors.New("consumer is nil")
	ErrConsumerPersonalDataRequired   = errors.New("consumer personal data is required")
	ErrConsumerContactDataRequired    = errors.New("consumer contact information is required")
	ErrConsumerCreditProfileRequired  = errors.New("consumer credit profile is required")
	ErrConsumerPrimaryAddressRequired = errors.New("consumer primary address id is required")
	ErrConsumerProductAlreadyLinked   = errors.New("product already contracted by consumer")
	ErrConsumerProductNotLinked       = errors.New("product not linked to consumer")
)

type Consumer struct {
	ID                   primitive.ObjectID
	Type                 valueobjects.ConsumerType
	PersonalData         ConsumerPersonalData
	CreditProfile        ConsumerCreditProfile
	Contact              ConsumerContactInformation
	PrimaryAddressID     primitive.ObjectID
	AdditionalAddressIDs []primitive.ObjectID
	ContractedProducts   []primitive.ObjectID
	UserID               primitive.ObjectID
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (c *Consumer) Validate() error {
	if c == nil {
		return ErrConsumerNil
	}

	if err := c.Type.Validate(); err != nil {
		return err
	}

	if err := c.PersonalData.Validate(c.Type); err != nil {
		return err
	}

	if err := c.Contact.Validate(); err != nil {
		return err
	}

	if err := c.CreditProfile.Validate(c.Type); err != nil {
		return err
	}

	if c.PrimaryAddressID.IsZero() {
		return ErrConsumerPrimaryAddressRequired
	}

	return nil
}

func (c *Consumer) HasContractedProduct(productID primitive.ObjectID) bool {
	if c == nil || productID.IsZero() {
		return false
	}

	for _, id := range c.ContractedProducts {
		if id == productID {
			return true
		}
	}

	return false
}

// HasLinkedUser reports whether the consumer already has an associated authentication profile.
func (c *Consumer) HasLinkedUser() bool {
	if c == nil {
		return false
	}
	return !c.UserID.IsZero()
}

func (c *Consumer) AddContractedProduct(productID primitive.ObjectID) error {
	if c == nil {
		return ErrConsumerNil
	}
	if productID.IsZero() {
		return fmt.Errorf("product id is required")
	}
	if c.HasContractedProduct(productID) {
		return ErrConsumerProductAlreadyLinked
	}

	c.ContractedProducts = append(c.ContractedProducts, productID)
	return nil
}

func (c *Consumer) RemoveContractedProduct(productID primitive.ObjectID) error {
	if c == nil {
		return ErrConsumerNil
	}
	if productID.IsZero() {
		return fmt.Errorf("product id is required")
	}
	if !c.HasContractedProduct(productID) {
		return ErrConsumerProductNotLinked
	}

	filtered := make([]primitive.ObjectID, 0, len(c.ContractedProducts))
	for _, id := range c.ContractedProducts {
		if id != productID {
			filtered = append(filtered, id)
		}
	}

	c.ContractedProducts = filtered
	return nil
}

// ConsumerPersonalData keeps personal/registration data for any consumer type.
type ConsumerPersonalData struct {
	Individual *ConsumerIndividualData
	Business   *ConsumerBusinessData
}

func (pd *ConsumerPersonalData) Validate(consumerType valueobjects.ConsumerType) error {
	if pd == nil {
		return ErrConsumerPersonalDataRequired
	}

	switch consumerType {
	case valueobjects.ConsumerTypeIndividual:
		if pd.Individual == nil {
			return ErrConsumerPersonalDataRequired
		}
		return pd.Individual.Validate()
	case valueobjects.ConsumerTypeBusiness:
		if pd.Business == nil {
			return ErrConsumerPersonalDataRequired
		}
		return pd.Business.Validate()
	default:
		return fmt.Errorf("unsupported consumer type: %s", consumerType)
	}
}

// ConsumerIndividualData represents required personal details for individuals.
type ConsumerIndividualData struct {
	FullName       string
	SocialName     string
	DocumentNumber string
	BirthDate      time.Time
	Nationality    string
	MaritalStatus  string
	Occupation     string
}

func (id ConsumerIndividualData) Validate() error {
	fullName := strings.TrimSpace(id.FullName)
	if fullName == "" {
		return fmt.Errorf("individual full name is required")
	}

	doc := strings.TrimSpace(id.DocumentNumber)
	if doc == "" {
		return fmt.Errorf("individual document number (CPF) is required")
	}
	if digits := countDigits(doc); digits != 11 {
		return fmt.Errorf("individual document number must contain 11 digits")
	}

	if id.BirthDate.IsZero() {
		return fmt.Errorf("individual birth date is required")
	}

	if time.Since(id.BirthDate) <= 0 {
		return fmt.Errorf("individual birth date must be in the past")
	}

	return nil
}

// ConsumerBusinessData represents required registry data for legal entities.
type ConsumerBusinessData struct {
	CorporateName     string
	TradeName         string
	DocumentNumber    string
	IncorporationDate time.Time
	LegalNature       string
	StateRegistration string
	MunicipalRegistry string
}

func (bd ConsumerBusinessData) Validate() error {
	corporateName := strings.TrimSpace(bd.CorporateName)
	if corporateName == "" {
		return fmt.Errorf("business corporate name is required")
	}

	doc := strings.TrimSpace(bd.DocumentNumber)
	if doc == "" {
		return fmt.Errorf("business document number (CNPJ) is required")
	}
	if digits := countDigits(doc); digits != 14 {
		return fmt.Errorf("business document number must contain 14 digits")
	}

	if bd.IncorporationDate.IsZero() {
		return fmt.Errorf("business incorporation date is required")
	}
	if time.Since(bd.IncorporationDate) <= 0 {
		return fmt.Errorf("business incorporation date must be in the past")
	}

	return nil
}

// ConsumerContactInformation stores contact channels used through the contracted financial journey.
type ConsumerContactInformation struct {
	Email          string
	Phone          string
	SecondaryPhone string
}

func (ci ConsumerContactInformation) Validate() error {
	email := strings.TrimSpace(ci.Email)
	if email == "" {
		return ErrConsumerContactDataRequired
	}
	if !strings.Contains(email, "@") {
		return fmt.Errorf("contact email appears to be invalid")
	}

	phone := strings.TrimSpace(ci.Phone)
	if phone == "" {
		return fmt.Errorf("primary contact phone is required")
	}

	return nil
}

// ConsumerCreditProfile aggregates financial indicators used for risk assessment and approvals.
type ConsumerCreditProfile struct {
	CreditScore             int
	MonthlyIncome           float64
	AnnualRevenue           float64
	CreditLimitRequested    float64
	CreditLimitApproved     float64
	OutstandingDebt         float64
	RiskLevel               string
	DelinquencyProbability  float64
	EmploymentStatus        string
	YearsInCurrentJob       int
	YearsInBusiness         int
	BankingRelationshipRank string
}

func (cp ConsumerCreditProfile) Validate(consumerType valueobjects.ConsumerType) error {
	if cp == (ConsumerCreditProfile{}) {
		return ErrConsumerCreditProfileRequired
	}

	if cp.CreditScore < 0 || cp.CreditScore > 1000 {
		return fmt.Errorf("credit score must be between 0 and 1000")
	}
	if cp.MonthlyIncome < 0 {
		return fmt.Errorf("monthly income cannot be negative")
	}
	if cp.AnnualRevenue < 0 {
		return fmt.Errorf("annual revenue cannot be negative")
	}
	if cp.CreditLimitRequested < 0 {
		return fmt.Errorf("requested credit limit cannot be negative")
	}
	if cp.CreditLimitApproved < 0 {
		return fmt.Errorf("approved credit limit cannot be negative")
	}
	if cp.OutstandingDebt < 0 {
		return fmt.Errorf("outstanding debt cannot be negative")
	}
	if cp.DelinquencyProbability < 0 || cp.DelinquencyProbability > 1 {
		return fmt.Errorf("delinquency probability must be between 0 and 1")
	}
	if cp.YearsInCurrentJob < 0 {
		return fmt.Errorf("years in current job cannot be negative")
	}
	if cp.YearsInBusiness < 0 {
		return fmt.Errorf("years in business cannot be negative")
	}

	if consumerType == valueobjects.ConsumerTypeIndividual {
		if cp.MonthlyIncome == 0 {
			return fmt.Errorf("monthly income is required for individual consumers")
		}
	}

	if consumerType == valueobjects.ConsumerTypeBusiness {
		if cp.AnnualRevenue == 0 {
			return fmt.Errorf("annual revenue is required for business consumers")
		}
	}

	return nil
}

func countDigits(value string) int {
	count := 0
	for _, r := range value {
		if r >= '0' && r <= '9' {
			count++
		}
	}
	return count
}
