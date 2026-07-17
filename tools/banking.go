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

func UpdateCardLimitToolDefinition() mcp.ToolDefinition {
	return mcp.ToolDefinition{
		Name:        "update_card_limit",
		Description: "Atualiza o limite do cartao de um cliente pelo customer_id.",
		Parameters: []mcp.ToolParameter{
			{
				Name:        "customer_id",
				Type:        "string",
				Description: "ID do cliente que tera o limite atualizado.",
				Required:    true,
			},
			{
				Name:        "new_limit",
				Type:        "integer",
				Description: "Novo limite total do cartao.",
				Required:    true,
			},
		},
	}
}

func CreatePixToolDefinition() mcp.ToolDefinition {
	return mcp.ToolDefinition{
		Name:        "create_pix",
		Description: "Cria uma transferencia PIX a partir da conta de um cliente. Operacao critica que exige confirmacao explicita do usuario.",
		Parameters: []mcp.ToolParameter{
			{
				Name:        "customer_id",
				Type:        "string",
				Description: "ID do cliente de origem do PIX.",
				Required:    true,
			},
			{
				Name:        "pix_key",
				Type:        "string",
				Description: "Chave PIX de destino.",
				Required:    true,
			},
			{
				Name:        "amount_cents",
				Type:        "integer",
				Description: "Valor do PIX em centavos.",
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

func UpdateCardLimitTool(user policy.AuthenticatedUser, customerID string, newLimit int) (banking.CardLimit, error) {

	decision := policy.CanUpdateCardLimit(user, customerID)

	if !decision.Allowed {
		return banking.CardLimit{}, errors.New(decision.Reason)
	}

	limit, err := banking.UpdateCardLimit(customerID, newLimit)
	if err != nil {
		return banking.CardLimit{}, err
	}

	return limit, nil
}

func CreatePixTool(user policy.AuthenticatedUser, customerID string, pixKey string, amountCents int, confirmed bool) (banking.PixResult, error) {
	if !confirmed {
		return banking.PixResult{}, errors.New("pix_requires_confirmation")
	}

	decision := policy.CanCreatePix(user, customerID)

	if !decision.Allowed {
		return banking.PixResult{}, errors.New(decision.Reason)
	}

	transaction, err := banking.CreatePix(customerID, pixKey, amountCents)
	if err != nil {
		return banking.PixResult{}, err
	}

	return transaction, nil
}
