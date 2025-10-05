package dto

import (
	"fmt"
	"strings"
	"time"

	"katseye/internal/domain/entities"
	valueobjects "katseye/internal/domain/value_objects"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	isoDateLayout     = "2006-01-02"
	isoDateTimeLayout = time.RFC3339
)

type ConsumerRequest struct {
	Type                 string                       `json:"type"`
	PersonalData         ConsumerPersonalDataRequest  `json:"personal_data"`
	CreditProfile        ConsumerCreditProfileRequest `json:"credit_profile"`
	Contact              ConsumerContactRequest       `json:"contact"`
	PrimaryAddressID     string                       `json:"primary_address_id"`
	AdditionalAddressIDs []string                     `json:"additional_address_ids"`
	ContractedProductIDs []string                     `json:"contracted_product_ids"`
	UserID               string                       `json:"user_id,omitempty"`
}

type ConsumerPersonalDataRequest struct {
	Individual *ConsumerIndividualDataRequest `json:"individual,omitempty"`
	Business   *ConsumerBusinessDataRequest   `json:"business,omitempty"`
}

type ConsumerIndividualDataRequest struct {
	FullName       string `json:"full_name"`
	SocialName     string `json:"social_name"`
	DocumentNumber string `json:"document_number"`
	BirthDate      string `json:"birth_date"`
	Nationality    string `json:"nationality"`
	MaritalStatus  string `json:"marital_status"`
	Occupation     string `json:"occupation"`
}

type ConsumerBusinessDataRequest struct {
	CorporateName     string `json:"corporate_name"`
	TradeName         string `json:"trade_name"`
	DocumentNumber    string `json:"document_number"`
	IncorporationDate string `json:"incorporation_date"`
	LegalNature       string `json:"legal_nature"`
	StateRegistration string `json:"state_registration"`
	MunicipalRegistry string `json:"municipal_registry"`
}

type ConsumerContactRequest struct {
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	SecondaryPhone string `json:"secondary_phone"`
}

type ConsumerCreditProfileRequest struct {
	CreditScore             int     `json:"credit_score"`
	MonthlyIncome           float64 `json:"monthly_income"`
	AnnualRevenue           float64 `json:"annual_revenue"`
	CreditLimitRequested    float64 `json:"credit_limit_requested"`
	CreditLimitApproved     float64 `json:"credit_limit_approved"`
	OutstandingDebt         float64 `json:"outstanding_debt"`
	RiskLevel               string  `json:"risk_level"`
	DelinquencyProbability  float64 `json:"delinquency_probability"`
	EmploymentStatus        string  `json:"employment_status"`
	YearsInCurrentJob       int     `json:"years_in_current_job"`
	YearsInBusiness         int     `json:"years_in_business"`
	BankingRelationshipRank string  `json:"banking_relationship_rank"`
}

type ConsumerResponse struct {
	ID                   string                        `json:"id"`
	Type                 string                        `json:"type"`
	PersonalData         ConsumerPersonalDataResponse  `json:"personal_data"`
	CreditProfile        ConsumerCreditProfileResponse `json:"credit_profile"`
	Contact              ConsumerContactResponse       `json:"contact"`
	PrimaryAddressID     string                        `json:"primary_address_id"`
	AdditionalAddressIDs []string                      `json:"additional_address_ids"`
	ContractedProducts   []string                      `json:"contracted_products"`
	UserID               string                        `json:"user_id,omitempty"`
	CreatedAt            time.Time                     `json:"created_at"`
	UpdatedAt            time.Time                     `json:"updated_at"`
}

type ConsumerPersonalDataResponse struct {
	Individual *ConsumerIndividualDataResponse `json:"individual,omitempty"`
	Business   *ConsumerBusinessDataResponse   `json:"business,omitempty"`
}

type ConsumerIndividualDataResponse struct {
	FullName       string    `json:"full_name"`
	SocialName     string    `json:"social_name,omitempty"`
	DocumentNumber string    `json:"document_number"`
	BirthDate      time.Time `json:"birth_date"`
	Nationality    string    `json:"nationality,omitempty"`
	MaritalStatus  string    `json:"marital_status,omitempty"`
	Occupation     string    `json:"occupation,omitempty"`
}

type ConsumerBusinessDataResponse struct {
	CorporateName     string    `json:"corporate_name"`
	TradeName         string    `json:"trade_name,omitempty"`
	DocumentNumber    string    `json:"document_number"`
	IncorporationDate time.Time `json:"incorporation_date"`
	LegalNature       string    `json:"legal_nature,omitempty"`
	StateRegistration string    `json:"state_registration,omitempty"`
	MunicipalRegistry string    `json:"municipal_registry,omitempty"`
}

type ConsumerContactResponse struct {
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	SecondaryPhone string `json:"secondary_phone,omitempty"`
}

type ConsumerCreditProfileResponse struct {
	CreditScore             int     `json:"credit_score"`
	MonthlyIncome           float64 `json:"monthly_income"`
	AnnualRevenue           float64 `json:"annual_revenue"`
	CreditLimitRequested    float64 `json:"credit_limit_requested"`
	CreditLimitApproved     float64 `json:"credit_limit_approved"`
	OutstandingDebt         float64 `json:"outstanding_debt"`
	RiskLevel               string  `json:"risk_level,omitempty"`
	DelinquencyProbability  float64 `json:"delinquency_probability"`
	EmploymentStatus        string  `json:"employment_status,omitempty"`
	YearsInCurrentJob       int     `json:"years_in_current_job"`
	YearsInBusiness         int     `json:"years_in_business"`
	BankingRelationshipRank string  `json:"banking_relationship_rank,omitempty"`
}

func (req *ConsumerRequest) ToEntity(id primitive.ObjectID) (*entities.Consumer, error) {
	if req == nil {
		return nil, fmt.Errorf("consumer request is nil")
	}

	consumerType, err := valueobjects.NewConsumerType(req.Type)
	if err != nil {
		return nil, fmt.Errorf("invalid consumer type: %w", err)
	}

	if strings.TrimSpace(req.PrimaryAddressID) == "" {
		return nil, fmt.Errorf("primary_address_id is required")
	}

	primaryAddressID, err := primitive.ObjectIDFromHex(req.PrimaryAddressID)
	if err != nil {
		return nil, fmt.Errorf("invalid primary address id: %w", err)
	}

	additionalAddressIDs := make([]primitive.ObjectID, 0, len(req.AdditionalAddressIDs))
	for _, idStr := range req.AdditionalAddressIDs {
		trimmed := strings.TrimSpace(idStr)
		if trimmed == "" {
			continue
		}
		addrID, addrErr := primitive.ObjectIDFromHex(trimmed)
		if addrErr != nil {
			return nil, fmt.Errorf("invalid additional address id: %w", addrErr)
		}
		additionalAddressIDs = append(additionalAddressIDs, addrID)
	}

	contractedProducts := make([]primitive.ObjectID, 0, len(req.ContractedProductIDs))
	for _, productID := range req.ContractedProductIDs {
		trimmed := strings.TrimSpace(productID)
		if trimmed == "" {
			continue
		}
		pid, pidErr := primitive.ObjectIDFromHex(trimmed)
		if pidErr != nil {
			return nil, fmt.Errorf("invalid contracted product id: %w", pidErr)
		}
		contractedProducts = append(contractedProducts, pid)
	}

	var userID primitive.ObjectID
	if trimmed := strings.TrimSpace(req.UserID); trimmed != "" {
		parsed, parseErr := primitive.ObjectIDFromHex(trimmed)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid user id: %w", parseErr)
		}
		userID = parsed
	}

	personalData, err := req.PersonalData.toEntity(consumerType)
	if err != nil {
		return nil, err
	}

	creditProfile, err := req.CreditProfile.toEntity(consumerType)
	if err != nil {
		return nil, err
	}

	contact := req.Contact.toEntity()

	consumer := &entities.Consumer{
		ID:                   id,
		Type:                 consumerType,
		PersonalData:         personalData,
		CreditProfile:        creditProfile,
		Contact:              contact,
		PrimaryAddressID:     primaryAddressID,
		AdditionalAddressIDs: additionalAddressIDs,
		ContractedProducts:   contractedProducts,
		UserID:               userID,
	}

	return consumer, nil
}

func (req ConsumerPersonalDataRequest) toEntity(consumerType valueobjects.ConsumerType) (entities.ConsumerPersonalData, error) {
	data := entities.ConsumerPersonalData{}

	switch consumerType {
	case valueobjects.ConsumerTypeIndividual:
		if req.Individual == nil {
			return data, fmt.Errorf("individual personal data is required for type individual")
		}
		birthDate, err := parseFlexibleDate(req.Individual.BirthDate)
		if err != nil {
			return data, fmt.Errorf("invalid birth date: %w", err)
		}
		data.Individual = &entities.ConsumerIndividualData{
			FullName:       req.Individual.FullName,
			SocialName:     req.Individual.SocialName,
			DocumentNumber: req.Individual.DocumentNumber,
			BirthDate:      birthDate,
			Nationality:    req.Individual.Nationality,
			MaritalStatus:  req.Individual.MaritalStatus,
			Occupation:     req.Individual.Occupation,
		}
	case valueobjects.ConsumerTypeBusiness:
		if req.Business == nil {
			return data, fmt.Errorf("business personal data is required for type business")
		}
		incorporationDate, err := parseFlexibleDate(req.Business.IncorporationDate)
		if err != nil {
			return data, fmt.Errorf("invalid incorporation date: %w", err)
		}
		data.Business = &entities.ConsumerBusinessData{
			CorporateName:     req.Business.CorporateName,
			TradeName:         req.Business.TradeName,
			DocumentNumber:    req.Business.DocumentNumber,
			IncorporationDate: incorporationDate,
			LegalNature:       req.Business.LegalNature,
			StateRegistration: req.Business.StateRegistration,
			MunicipalRegistry: req.Business.MunicipalRegistry,
		}
	default:
		return data, fmt.Errorf("unsupported consumer type: %s", consumerType)
	}

	return data, nil
}

func (req ConsumerCreditProfileRequest) toEntity(consumerType valueobjects.ConsumerType) (entities.ConsumerCreditProfile, error) {
	profile := entities.ConsumerCreditProfile{
		CreditScore:             req.CreditScore,
		MonthlyIncome:           req.MonthlyIncome,
		AnnualRevenue:           req.AnnualRevenue,
		CreditLimitRequested:    req.CreditLimitRequested,
		CreditLimitApproved:     req.CreditLimitApproved,
		OutstandingDebt:         req.OutstandingDebt,
		RiskLevel:               req.RiskLevel,
		DelinquencyProbability:  req.DelinquencyProbability,
		EmploymentStatus:        req.EmploymentStatus,
		YearsInCurrentJob:       req.YearsInCurrentJob,
		YearsInBusiness:         req.YearsInBusiness,
		BankingRelationshipRank: req.BankingRelationshipRank,
	}

	return profile, nil
}

func (req ConsumerContactRequest) toEntity() entities.ConsumerContactInformation {
	return entities.ConsumerContactInformation{
		Email:          req.Email,
		Phone:          req.Phone,
		SecondaryPhone: req.SecondaryPhone,
	}
}

func NewConsumerResponse(consumer *entities.Consumer) ConsumerResponse {
	if consumer == nil {
		return ConsumerResponse{}
	}

	response := ConsumerResponse{
		ID:                   consumer.ID.Hex(),
		Type:                 consumer.Type.String(),
		PersonalData:         newPersonalDataResponse(consumer.PersonalData),
		CreditProfile:        newCreditProfileResponse(consumer.CreditProfile),
		Contact:              newContactResponse(consumer.Contact),
		PrimaryAddressID:     consumer.PrimaryAddressID.Hex(),
		AdditionalAddressIDs: objectIDSliceToHex(consumer.AdditionalAddressIDs),
		ContractedProducts:   objectIDSliceToHex(consumer.ContractedProducts),
		CreatedAt:            consumer.CreatedAt,
		UpdatedAt:            consumer.UpdatedAt,
	}

	if !consumer.UserID.IsZero() {
		response.UserID = consumer.UserID.Hex()
	}

	return response
}

func NewConsumerResponseList(consumers []*entities.Consumer) []ConsumerResponse {
	if len(consumers) == 0 {
		return nil
	}

	responses := make([]ConsumerResponse, 0, len(consumers))
	for _, consumer := range consumers {
		responses = append(responses, NewConsumerResponse(consumer))
	}

	return responses
}

func newPersonalDataResponse(data entities.ConsumerPersonalData) ConsumerPersonalDataResponse {
	response := ConsumerPersonalDataResponse{}

	if data.Individual != nil {
		individual := data.Individual
		response.Individual = &ConsumerIndividualDataResponse{
			FullName:       individual.FullName,
			SocialName:     individual.SocialName,
			DocumentNumber: individual.DocumentNumber,
			BirthDate:      individual.BirthDate,
			Nationality:    individual.Nationality,
			MaritalStatus:  individual.MaritalStatus,
			Occupation:     individual.Occupation,
		}
	}

	if data.Business != nil {
		business := data.Business
		response.Business = &ConsumerBusinessDataResponse{
			CorporateName:     business.CorporateName,
			TradeName:         business.TradeName,
			DocumentNumber:    business.DocumentNumber,
			IncorporationDate: business.IncorporationDate,
			LegalNature:       business.LegalNature,
			StateRegistration: business.StateRegistration,
			MunicipalRegistry: business.MunicipalRegistry,
		}
	}

	return response
}

func newContactResponse(contact entities.ConsumerContactInformation) ConsumerContactResponse {
	return ConsumerContactResponse{
		Email:          contact.Email,
		Phone:          contact.Phone,
		SecondaryPhone: contact.SecondaryPhone,
	}
}

func newCreditProfileResponse(profile entities.ConsumerCreditProfile) ConsumerCreditProfileResponse {
	return ConsumerCreditProfileResponse{
		CreditScore:             profile.CreditScore,
		MonthlyIncome:           profile.MonthlyIncome,
		AnnualRevenue:           profile.AnnualRevenue,
		CreditLimitRequested:    profile.CreditLimitRequested,
		CreditLimitApproved:     profile.CreditLimitApproved,
		OutstandingDebt:         profile.OutstandingDebt,
		RiskLevel:               profile.RiskLevel,
		DelinquencyProbability:  profile.DelinquencyProbability,
		EmploymentStatus:        profile.EmploymentStatus,
		YearsInCurrentJob:       profile.YearsInCurrentJob,
		YearsInBusiness:         profile.YearsInBusiness,
		BankingRelationshipRank: profile.BankingRelationshipRank,
	}
}

func objectIDSliceToHex(values []primitive.ObjectID) []string {
	if len(values) == 0 {
		return nil
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value.IsZero() {
			continue
		}
		result = append(result, value.Hex())
	}
	return result
}

func parseFlexibleDate(value string) (time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return time.Time{}, fmt.Errorf("empty date value")
	}

	layouts := []string{isoDateTimeLayout, isoDateLayout}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, trimmed); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date format, expected YYYY-MM-DD or RFC3339")
}
