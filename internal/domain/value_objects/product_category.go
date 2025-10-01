package valueobjects

type ProductCategory string

const (
	ProductCategoryPersonal ProductCategory = "personal"
	ProductCategoryBusiness ProductCategory = "business"
	ProductCategoryOthers   ProductCategory = "others"
)

// internal/domain/value_objects/risk_level.go
type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "low"
	RiskLevelMedium RiskLevel = "medium"
	RiskLevelHigh   RiskLevel = "high"
)
