package valueobjects

import (
	"errors"
	"strings"
)

type ConsumerType string

const (
	ConsumerTypeIndividual ConsumerType = "individual"
	ConsumerTypeBusiness   ConsumerType = "business"
)

var (
	ErrInvalidConsumerType = errors.New("invalid consumer type")

	consumerTypeAliases = map[string]ConsumerType{
		"individual":     ConsumerTypeIndividual,
		"pf":             ConsumerTypeIndividual,
		"pessoa_fisica":  ConsumerTypeIndividual,
		"fisica":         ConsumerTypeIndividual,
		"natural_person": ConsumerTypeIndividual,

		"business":        ConsumerTypeBusiness,
		"pj":              ConsumerTypeBusiness,
		"pessoa_juridica": ConsumerTypeBusiness,
		"juridica":        ConsumerTypeBusiness,
		"corporate":       ConsumerTypeBusiness,
	}
)

func NewConsumerType(value string) (ConsumerType, error) {
	normalized := strings.TrimSpace(strings.ToLower(value))
	if normalized == "" {
		return "", ErrInvalidConsumerType
	}

	if mapped, ok := consumerTypeAliases[normalized]; ok {
		return mapped, nil
	}

	consumerType := ConsumerType(normalized)
	if err := consumerType.Validate(); err != nil {
		return "", err
	}

	return consumerType, nil
}

func (ct ConsumerType) Validate() error {
	switch ct {
	case ConsumerTypeIndividual, ConsumerTypeBusiness:
		return nil
	default:
		return ErrInvalidConsumerType
	}
}

func (ct ConsumerType) String() string {
	return string(ct)
}

func (ct ConsumerType) IsIndividual() bool {
	return ct == ConsumerTypeIndividual
}

func (ct ConsumerType) IsBusiness() bool {
	return ct == ConsumerTypeBusiness
}
