package models

import (
	"time"

	"katseye/internal/domain/entities"
	valueobjects "katseye/internal/domain/value_objects"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsumerDocument struct {
	ID                   primitive.ObjectID            `bson:"_id,omitempty"`
	Type                 valueobjects.ConsumerType     `bson:"type"`
	PersonalData         ConsumerPersonalDataDocument  `bson:"personal_data"`
	CreditProfile        ConsumerCreditProfileDocument `bson:"credit_profile"`
	Contact              ConsumerContactDocument       `bson:"contact"`
	PrimaryAddressID     primitive.ObjectID            `bson:"primary_address_id"`
	AdditionalAddressIDs []primitive.ObjectID          `bson:"additional_address_ids,omitempty"`
	ContractedProducts   []primitive.ObjectID          `bson:"contracted_products,omitempty"`
	CreatedAt            time.Time                     `bson:"created_at"`
	UpdatedAt            time.Time                     `bson:"updated_at"`
}

type ConsumerPersonalDataDocument struct {
	Individual *ConsumerIndividualDataDocument `bson:"individual,omitempty"`
	Business   *ConsumerBusinessDataDocument   `bson:"business,omitempty"`
}

type ConsumerIndividualDataDocument struct {
	FullName       string    `bson:"full_name"`
	SocialName     string    `bson:"social_name,omitempty"`
	DocumentNumber string    `bson:"document_number"`
	BirthDate      time.Time `bson:"birth_date"`
	Nationality    string    `bson:"nationality,omitempty"`
	MaritalStatus  string    `bson:"marital_status,omitempty"`
	Occupation     string    `bson:"occupation,omitempty"`
}

type ConsumerBusinessDataDocument struct {
	CorporateName     string    `bson:"corporate_name"`
	TradeName         string    `bson:"trade_name,omitempty"`
	DocumentNumber    string    `bson:"document_number"`
	IncorporationDate time.Time `bson:"incorporation_date"`
	LegalNature       string    `bson:"legal_nature,omitempty"`
	StateRegistration string    `bson:"state_registration,omitempty"`
	MunicipalRegistry string    `bson:"municipal_registry,omitempty"`
}

type ConsumerContactDocument struct {
	Email          string `bson:"email"`
	Phone          string `bson:"phone"`
	SecondaryPhone string `bson:"secondary_phone,omitempty"`
}

type ConsumerCreditProfileDocument struct {
	CreditScore             int     `bson:"credit_score"`
	MonthlyIncome           float64 `bson:"monthly_income"`
	AnnualRevenue           float64 `bson:"annual_revenue"`
	CreditLimitRequested    float64 `bson:"credit_limit_requested"`
	CreditLimitApproved     float64 `bson:"credit_limit_approved"`
	OutstandingDebt         float64 `bson:"outstanding_debt"`
	RiskLevel               string  `bson:"risk_level,omitempty"`
	DelinquencyProbability  float64 `bson:"delinquency_probability"`
	EmploymentStatus        string  `bson:"employment_status,omitempty"`
	YearsInCurrentJob       int     `bson:"years_in_current_job"`
	YearsInBusiness         int     `bson:"years_in_business"`
	BankingRelationshipRank string  `bson:"banking_relationship_rank,omitempty"`
}

func NewConsumerDocument(consumer *entities.Consumer) ConsumerDocument {
	if consumer == nil {
		return ConsumerDocument{}
	}

	doc := ConsumerDocument{
		ID:                   consumer.ID,
		Type:                 consumer.Type,
		PersonalData:         newPersonalDataDocument(consumer.PersonalData),
		CreditProfile:        newCreditProfileDocument(consumer.CreditProfile),
		Contact:              newContactDocument(consumer.Contact),
		PrimaryAddressID:     consumer.PrimaryAddressID,
		AdditionalAddressIDs: append([]primitive.ObjectID(nil), consumer.AdditionalAddressIDs...),
		ContractedProducts:   append([]primitive.ObjectID(nil), consumer.ContractedProducts...),
		CreatedAt:            consumer.CreatedAt,
		UpdatedAt:            consumer.UpdatedAt,
	}

	return doc
}

func (doc ConsumerDocument) ToEntity() *entities.Consumer {
	consumer := &entities.Consumer{
		ID:                   doc.ID,
		Type:                 doc.Type,
		PersonalData:         doc.PersonalData.toEntity(),
		CreditProfile:        doc.CreditProfile.toEntity(),
		Contact:              doc.Contact.toEntity(),
		PrimaryAddressID:     doc.PrimaryAddressID,
		AdditionalAddressIDs: append([]primitive.ObjectID(nil), doc.AdditionalAddressIDs...),
		ContractedProducts:   append([]primitive.ObjectID(nil), doc.ContractedProducts...),
		CreatedAt:            doc.CreatedAt,
		UpdatedAt:            doc.UpdatedAt,
	}

	return consumer
}

func newPersonalDataDocument(data entities.ConsumerPersonalData) ConsumerPersonalDataDocument {
	doc := ConsumerPersonalDataDocument{}

	if data.Individual != nil {
		individual := *data.Individual
		doc.Individual = &ConsumerIndividualDataDocument{
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
		business := *data.Business
		doc.Business = &ConsumerBusinessDataDocument{
			CorporateName:     business.CorporateName,
			TradeName:         business.TradeName,
			DocumentNumber:    business.DocumentNumber,
			IncorporationDate: business.IncorporationDate,
			LegalNature:       business.LegalNature,
			StateRegistration: business.StateRegistration,
			MunicipalRegistry: business.MunicipalRegistry,
		}
	}

	return doc
}

func (doc ConsumerPersonalDataDocument) toEntity() entities.ConsumerPersonalData {
	personalData := entities.ConsumerPersonalData{}

	if doc.Individual != nil {
		individual := entities.ConsumerIndividualData{
			FullName:       doc.Individual.FullName,
			SocialName:     doc.Individual.SocialName,
			DocumentNumber: doc.Individual.DocumentNumber,
			BirthDate:      doc.Individual.BirthDate,
			Nationality:    doc.Individual.Nationality,
			MaritalStatus:  doc.Individual.MaritalStatus,
			Occupation:     doc.Individual.Occupation,
		}
		personalData.Individual = &individual
	}

	if doc.Business != nil {
		business := entities.ConsumerBusinessData{
			CorporateName:     doc.Business.CorporateName,
			TradeName:         doc.Business.TradeName,
			DocumentNumber:    doc.Business.DocumentNumber,
			IncorporationDate: doc.Business.IncorporationDate,
			LegalNature:       doc.Business.LegalNature,
			StateRegistration: doc.Business.StateRegistration,
			MunicipalRegistry: doc.Business.MunicipalRegistry,
		}
		personalData.Business = &business
	}

	return personalData
}

func newContactDocument(contact entities.ConsumerContactInformation) ConsumerContactDocument {
	return ConsumerContactDocument{
		Email:          contact.Email,
		Phone:          contact.Phone,
		SecondaryPhone: contact.SecondaryPhone,
	}
}

func (doc ConsumerContactDocument) toEntity() entities.ConsumerContactInformation {
	return entities.ConsumerContactInformation{
		Email:          doc.Email,
		Phone:          doc.Phone,
		SecondaryPhone: doc.SecondaryPhone,
	}
}

func newCreditProfileDocument(profile entities.ConsumerCreditProfile) ConsumerCreditProfileDocument {
	return ConsumerCreditProfileDocument{
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

func (doc ConsumerCreditProfileDocument) toEntity() entities.ConsumerCreditProfile {
	return entities.ConsumerCreditProfile{
		CreditScore:             doc.CreditScore,
		MonthlyIncome:           doc.MonthlyIncome,
		AnnualRevenue:           doc.AnnualRevenue,
		CreditLimitRequested:    doc.CreditLimitRequested,
		CreditLimitApproved:     doc.CreditLimitApproved,
		OutstandingDebt:         doc.OutstandingDebt,
		RiskLevel:               doc.RiskLevel,
		DelinquencyProbability:  doc.DelinquencyProbability,
		EmploymentStatus:        doc.EmploymentStatus,
		YearsInCurrentJob:       doc.YearsInCurrentJob,
		YearsInBusiness:         doc.YearsInBusiness,
		BankingRelationshipRank: doc.BankingRelationshipRank,
	}
}
