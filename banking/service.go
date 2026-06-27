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

func GetCustomerProfile(customerID string) (CustomerProfile, error) {

	profile, ok := customers[customerID]

	if !ok {
		return CustomerProfile{}, fmt.Errorf("cliente não encontrado: %s", customerID)
	}

	return profile, nil

}
