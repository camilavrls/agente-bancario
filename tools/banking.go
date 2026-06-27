package tools

import (
	"errors"

	"agente-bancario/banking"
	"agente-bancario/mcp"
	"agente-bancario/policy"
)

func GetCustomerProfileToolDefinition() mcp.ToolDefinition {
	return mcp.ToolDefinition{
		Name:        "get_customer_profile",
		Description: "Consulta o perfil de um cliente pelo customer_id, retornando segmento e score de credito.",
		Parameters: []mcp.ToolParameter{
			{
				Name:        "customer_id",
				Type:        "string",
				Description: "ID do cliente que deseja consultar.",
				Required:    true,
			},
		},
	}
}

func GetCustomerProfileTool(user policy.AuthenticatedUser, customerID string) (banking.CustomerProfile, error) {

	decision := policy.CanAccessCustomerProfile(user, customerID)

	if !decision.Allowed {
		return banking.CustomerProfile{}, errors.New(decision.Reason)
	}

	profile, err := banking.GetCustomerProfile(customerID)
	if err != nil {
		return banking.CustomerProfile{}, err
	}

	return profile, nil
}
