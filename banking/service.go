package banking

import "fmt"

var customers = map[string]CustomerProfile{
	"cust-123": {
		ID:          "cust-123",
		Name:        "Maria",
		Segment:     "Uniclass",
		CreditScore: 820,
	},

	"cust-456": {
		ID:          "cust-456",
		Name:        "Joao",
		Segment:     "Personnalité",
		CreditScore: 300,
	},
}

var cardLimits = map[string]CardLimit{
	"cust-123": {
		CustomerID:      "cust-123",
		CurrentLimit:    10000,
		AvailableLimit:  7400,
		MaxAllowedLimit: 15000,
	},

	"cust-456": {
		CustomerID:      "cust-456",
		CurrentLimit:    3000,
		AvailableLimit:  850,
		MaxAllowedLimit: 5000,
	},
}

func GetCustomerProfile(customerID string) (CustomerProfile, error) {
	profile, ok := customers[customerID]
	if !ok {
		return CustomerProfile{}, fmt.Errorf("cliente não encontrado: %s", customerID)
	}

	return profile, nil

}

func GetCardLimit(customerID string) (CardLimit, error) {
	limit, ok := cardLimits[customerID]
	if !ok {
		return CardLimit{}, fmt.Errorf("limite de cartão não encontrado: %s", customerID)
	}

	return limit, nil

}

func UpdateCardLimit(customerID string, newLimit int) (CardLimit, error) {
	limit, ok := cardLimits[customerID]
	if !ok {
		return CardLimit{}, fmt.Errorf("limite de cartão não encontrado: %s", customerID)
	}

	if newLimit <= 0 {
		return CardLimit{}, fmt.Errorf("novo limite deve ser maior que zero")
	}

	if newLimit > limit.MaxAllowedLimit {
		return CardLimit{}, fmt.Errorf("novo limite excede o limite maximo permitido")
	}

	usedLimit := limit.CurrentLimit - limit.AvailableLimit
	if newLimit < usedLimit {
		return CardLimit{}, fmt.Errorf("novo limite nao pode ser menor que o valor ja utilizado")
	}

	limit.CurrentLimit = newLimit
	limit.AvailableLimit = newLimit - usedLimit
	cardLimits[customerID] = limit

	return limit, nil
}
