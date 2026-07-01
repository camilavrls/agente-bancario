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

func GetCardLimitToolDefinition() mcp.ToolDefinition {
	return mcp.ToolDefinition{
		Name:        "get_card_limit",
		Description: "Consulta o limite atual e disponivel do cartao de um cliente pelo customer_id.",
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

	decision := policy.CanAccessCustomerResource(user, customerID)

	if !decision.Allowed {
		return banking.CustomerProfile{}, errors.New(decision.Reason)
	}

	profile, err := banking.GetCustomerProfile(customerID)
	if err != nil {
		return banking.CustomerProfile{}, err
	}

	return profile, nil
}

func GetCardLimitTool(user policy.AuthenticatedUser, customerID string) (banking.CardLimit, error) {

	decision := policy.CanAccessCustomerResource(user, customerID)

	if !decision.Allowed {
		return banking.CardLimit{}, errors.New(decision.Reason)
	}

	limit, err := banking.GetCardLimit(customerID)
	if err != nil {
		return banking.CardLimit{}, err
	}

	return limit, nil
}
